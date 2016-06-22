package core

import (
	"log"
	"time"

	"internal/pref"

	"github.com/nats-io/nats"
)

const (
	bucketStreamIn  = "stream-in"
	bucketStreamOut = "stream-out"
	bucketStreamErr = "stream-err"

	subjectSteamIn  = "m12." + bucketStreamIn
	subjectSteamOut = "m12." + bucketStreamOut

	listN = 100
	tickD = 10 * time.Second
)

func Init() error {
	var err error
	err = initREDIS(pref.Redis)
	if err != nil {
		return err
	}

	err = initMINIO(pref.Minio)
	if err != nil {
		return err
	}

	return initNATS(pref.NATS)
}

func initREDIS(addr string) error {
	var err error
	pREDIS, err = openREDIS(addr)
	if err != nil {
		return err
	}

	return waitDBFromDisk(1 * time.Second)
}

func waitDBFromDisk(d time.Duration) error {
	c := pREDIS.Get()
	defer c.Close()

	t := time.NewTicker(d)
	var err error
	for range t.C {
		_, err = c.Do("PING")
		if err != nil {
			log.Println(err)
			continue
		}
		break
	}
	t.Stop()
	return err
}

func initMINIO(addr string) error {
	var err error
	cMINIO, err = openMINIO(addr)
	if err != nil {
		return err
	}

	return makeBuckets(bucketStreamIn, bucketStreamOut, bucketStreamErr)
}

func makeBuckets(list ...string) error {
	var err error
	for i := range list {
		b := list[i]
		err = cMINIO.BucketExists(b)
		if err != nil {
			err = cMINIO.MakeBucket(b, "")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func initNATS(addr string) error {
	var err error

	cNATS, err = openNATS(addr)
	if err != nil {
		return err
	}

	_, err = cNATS.Subscribe(subjectSteamIn, func(m *nats.Msg) {
		proc(m.Data)
	})
	if err != nil {
		return err
	}

	go publishing(bucketStreamOut, subjectSteamOut, listN, tickD)
	go publishing(bucketStreamIn, subjectSteamIn, listN, tickD)

	return nil
}

func publishing(bucket, subject string, n int, d time.Duration) {
	var err error
	for range time.Tick(d) {
		err = publish(bucket, subject, n)
		if err != nil {
			log.Println(err)
		}
	}
}

func publish(bucket, subject string, n int) error {
	l, err := listObjectsN(bucket, n)
	if err != nil {
		return err
	}

	m := make([][]byte, len(l))
	for i := range l {
		log.Println("DEBUG", bucket, l[i])
		m[i] = pair{bucket, l[i]}.marshal()
	}

	return publishEach(subject, m...)
}

func publishEach(subject string, msgs ...[]byte) error {
	var err error
	for i := range msgs {
		err = cNATS.Publish(subject, msgs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func listObjectsN(bucket string, n int) ([]string, error) {
	doneCh := make(chan struct{}, 1)
	defer func() { close(doneCh) }()

	i := 0
	out := make([]string, 0, n)
	for o := range cMINIO.ListObjects(bucket, "", false, doneCh) {
		if o.Err != nil {
			return nil, o.Err
		}
		if len(out) < n { // workaround
			out = append(out, o.Key)
		}
		i++
		if i == n {
			doneCh <- struct{}{}
		}
	}

	return out, nil
}
