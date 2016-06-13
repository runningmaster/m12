package core

import (
	"log"
	"time"

	"internal/minio"
	"internal/pref"
	"internal/redis"

	"github.com/nats-io/nats"
)

func Init() error {
	var err error

	err = minio.Run(pref.Minio)
	if err != nil {
		return err
	}

	err = redis.Run(pref.Redis)
	if err != nil {
		return err
	}

	err = minio.InitBacketList(backetStreamIn, backetStreamOut, backetStreamErr)
	if err != nil {
		return err
	}
	return initNATS(pref.NATS)

}

func initNATS(addr string) error {
	var err error

	cNATS, err = openNATS(addr)
	if err != nil {
		return err
	}

	_, err = cNATS.Subscribe(backetStreamIn, func(m *nats.Msg) {
		err := proc(m.Data)
		if err != nil {
			log.Println("ListenAndServe", err)
			//goToStreamErr(m.ID, ?) // FIXME
		}
		// remove object
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
	l, err := minio.ListObjects(backet, n)
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
