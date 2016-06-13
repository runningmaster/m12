package nats

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"

	"github.com/nats-io/nats"
)

var cli *nats.Conn

func Run(addr string) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	var opts []nats.Option
	if u.User != nil {
		opts = append(opts, nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	}

	cli, err = nats.Connect(addr, opts)
	if err != nil {
		return fmt.Errorf("nats: %s", err)
	}

	return nil
}

func ListenAndServe(subject string, serveFunc func([]byte) error) error {
	_, err := cli.Subscribe(subject, func(m *nats.Msg) {
		if serveFunc == nil {
			return
		}
		err := serveFunc(m.Data)
		if err != nil {
			log.Println("ListenAndServe", err)
			//goToStreamErr(m.ID, ?) // FIXME
		}
	})
	return err
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
