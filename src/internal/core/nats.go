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
	testNATSConsumer()
	goNotifyStream(10)

	return nil
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
		log.Println(objs[i].Key)
		//natsCli.Publish(subject, []byte(objs[i].Key))
	}

	//natsCli.Publish(subject, []byte(strToSHA1(time.Now().String())))
	return nil
}

func testNATSConsumer() {
	natsCli.Subscribe(flag.NATSSubjectSteamIn, func(m *nats.Msg) {
		fmt.Printf("NATS save: %s %s %s\n", time.Now().String(), flag.NATSSubjectSteamIn, string(m.Data))
	})

	natsCli.Subscribe(flag.NATSSubjectSteamOut, func(m *nats.Msg) {
		fmt.Printf("NATS save: %s %s %s\n", time.Now().String(), flag.NATSSubjectSteamOut, string(m.Data))
	})
}
