package nats

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nats-io/nats"
)

var (
	cli    *nats.Conn
	logger = log.New(ioutil.Discard, "", log.LstdFlags)
)

func Run(addr string, l *log.Logger) error {
	if l != nil {
		logger = l
	}

	var err error
	cli, err = nats.Connect(addr, nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return fmt.Errorf("nats: %s", err)
	}

	return nil
}

func ListenAndServe(subject string, serveFunc func([]byte) error) {
	cli.Subscribe(subject, func(m *nats.Msg) {
		if serveFunc == nil {
			return
		}
		err := serveFunc(m.Data)
		if err != nil {
			logger.Println(err)
		}
	})
}

func PublishEach(subject string, list ...[]byte) error {
	var err error
	for i := range list {
		err = cli.Publish(subject, list[i])
		if err != nil {
			return err
		}
	}
	return nil
}
