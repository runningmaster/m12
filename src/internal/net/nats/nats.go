package nats

import (
	"crypto/tls"
	"log"
	"net/url"
	"time"

	nats "github.com/nats-io/go-nats"
)

var cli *nats.Conn

// Init inits client for NATS Server
func Init(addr string) error {
	c, err := makeConn(addr)
	if err != nil {
		return err
	}
	cli = c
	return nil
}

func makeConn(addr string) (*nats.Conn, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	opts := []nats.Option{nats.MaxReconnects(-1)}
	if u.User != nil {
		opts = append(opts, nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	}

	c, err := nats.Connect(addr, opts...)
	// workaround for system reboot
	if err != nil {
		log.Println(err)
		time.Sleep(10 * time.Second)
		c, err = nats.Connect(addr, opts...)
	}
	if err != nil {
		return nil, err
	}

	return c, nil
}

func Subscribe(s string, f func([]byte)) error {
	_, err := cli.Subscribe(s, func(m *nats.Msg) {
		go f(m.Data)
	})
	return err
}

func Publish(s string, data []byte) error {
	return cli.Publish(s, data)
}
