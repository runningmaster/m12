package core

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"internal/flag"

	"github.com/nats-io/nats"
)

const sendN = 10

var natsCli *nats.Conn

func initNATSCli() error {
	var err error

	natsCli, err = nats.Connect(flag.NATS, nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return fmt.Errorf("core: nats: %s", err)
	}

	goListenToNATS()
	goNotifyStream(sendN)

	return nil
}

func goListenToNATS() {
	natsCli.Subscribe(flag.NATSSubjectSteamIn, func(m *nats.Msg) {
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

	var b []byte
	for i := range objs {
		b, _ = pair{backet, objs[i].Key}.packToJSON()
		natsCli.Publish(subject, b)
	}

	return nil
}
