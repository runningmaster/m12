package core

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"internal/flag"

	"github.com/nats-io/nats"
)

var natsCli *nats.Conn

func initNATSCli() error {
	var err error

	natsCli, err = nats.Connect(flag.NATS, nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return fmt.Errorf("core: nats: %s", err)
	}

	goListenToNATS()
	goNotifyStream(10)

	return nil
}

func goListenToNATS() {
	natsCli.Subscribe(flag.NATSSubjectSteamIn, func(m *nats.Msg) {
		go proc(m.Data)
	})
}

func goNotifyStream(n int) {
	go func() {
		c := time.Tick(1 * time.Second)
		var err error
		for _ = range c {
			err = notifyStream(backetStreamIn, flag.NATSSubjectSteamIn, n)
			if err != nil {
				log.Println(err)
			}
			err = notifyStream(backetStreamOut, flag.NATSSubjectSteamOut, n)
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

func notifyStream(backet, subject string, n int) error {
	objs, err := listObjectsN(backet, "", false, n)
	if err != nil {
		return err
	}

	for i := range objs {
		natsCli.Publish(subject, pathS3{backet, objs[i].Key}.makeReadCloser())
	}

	return nil
}
