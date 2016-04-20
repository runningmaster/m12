package core

import (
	"crypto/tls"
	"fmt"
	"time"

	"internal/flag"

	"github.com/nats-io/nats"
)

var (
	natsCli *nats.Conn
)

func initNATSCli() error {
	var err error
	natsCli, err = nats.Connect(flag.NATS, nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return fmt.Errorf("core: nats: %s", err)
	}
	fmt.Println(natsCli.MaxPayload())
	//natsCli.Close()
	go testNATSProducer()
	testNATSConsumer()
	return nil
}

func testNATSProducer() {
	c := time.Tick(1 * time.Second)
	n := 0
	s := ""
	for _ = range c {
		n++
		s = fmt.Sprintf("Hello World! %d", n)
		fmt.Println("Send", s)
		natsCli.Publish("foo", []byte(s))

	}
}

func testNATSConsumer() {
	natsCli.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	//	fmt.Println("testNATSConsumer")
}
