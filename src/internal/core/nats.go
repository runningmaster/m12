package core

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats"
)

const (
	sendN = 10

	subjectSteamIn  = backetStreamIn + ".67a7ea16"
	subjectSteamOut = backetStreamOut + ".0566ce58"
)

var natsCli *nats.Conn

func initNATSCli(addr string) error {
	var err error

	natsCli, err = nats.Connect(addr, nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return fmt.Errorf("core: nats: %s", err)
	}

	goListenToNATS()
	goNotifyStream(sendN)

	return nil
}

func goListenToNATS() {
	natsCli.Subscribe(subjectSteamIn, func(m *nats.Msg) {
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
			err = notifyStream(backetStreamIn, subjectSteamIn, n)
			if err != nil {
				log.Println(err)
			}
			err = notifyStream(backetStreamOut, subjectSteamOut, n)
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
