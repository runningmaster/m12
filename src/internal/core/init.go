package core

import (
	"log"
	"time"

	"internal/pref"
	"internal/redis"

	"github.com/nats-io/nats"
)

func Init() error {
	var err error
	err = redis.Run(pref.Redis)
	if err != nil {
		return err
	}

	err = initMINIO(pref.Minio)
	if err != nil {
		return err
	}

	return initNATS(pref.NATS)
}

func initMINIO(addr string) error {
	var err error
	cMINIO, err = openMINIO(addr)
	if err != nil {
		return err
	}

	return makeBackets(backetStreamIn, backetStreamOut, backetStreamErr)
}

func makeBackets(list ...string) error {
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

	go publishing(backetStreamOut, subjectSteamOut, listN, tickD)
	go publishing(backetStreamIn, subjectSteamIn, listN, tickD)

	return nil
}

func publishing(backet, subject string, n int, d time.Duration) {
	var err error
	for range time.Tick(d) {
		err = publish(backet, subject, n)
		if err != nil {
			log.Println(err)
		}
	}
}

func publish(backet, subject string, n int) error {
	l, err := listObjectsN(backet, n)
	if err != nil {
		return err
	}

	m := make([][]byte, len(l))
	for i := range l {
		m[i] = pair{backet, l[i]}.marshalJSON()
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

func listObjectsN(backet string, n int) ([]string, error) {
	doneCh := make(chan struct{}, 1)
	defer func() { close(doneCh) }()

	i := 0
	out := make([]string, 0, n)
	for o := range cMINIO.ListObjects(backet, "", false, doneCh) {
		if o.Err != nil {
			return nil, o.Err
		}
		out = append(out, o.Key)
		i++
		if i == n {
			doneCh <- struct{}{}
		}
	}

	return out, nil
}
