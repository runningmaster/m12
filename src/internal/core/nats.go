package core

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"internal/conf"

	"github.com/nats-io/nats"
)

const sendN = 10

var natsCli *nats.Conn

func initNATSCli() error {
	var err error

	natsCli, err = nats.Connect(conf.NATS, nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return fmt.Errorf("core: nats: %s", err)
	}

	goListenToNATS()
	goNotifyStream(sendN)

	return nil
}

func goListenToNATS() {
	natsCli.Subscribe(conf.NATSSubjectSteamIn, func(m *nats.Msg) {
		go func() {
			p, err := makePairFromJSON(m.Data)
			if err != nil {
				log.Println(err)
			}
			err = proc(p.Backet, p.Object)
			if err != nil {
				log.Println(err)
			}
		}()
	})
}

func goNotifyStream(n int) {
	go func() {
		c := time.Tick(1 * time.Second)
		var err error
		for _ = range c {
			err = notifyStream(backetStreamIn, conf.NATSSubjectSteamIn, n)
			if err != nil {
				log.Println(err)
			}
			err = notifyStream(backetStreamOut, conf.NATSSubjectSteamOut, n)
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
		natsCli.Publish(subject, pair{backet, objs[i].Key}.packToJSON())
	}

	return nil
}
