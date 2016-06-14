package core

import (
	"crypto/tls"
	"net/url"

	minio "github.com/minio/minio-go"
	"github.com/nats-io/nats"
)

var cNATS *nats.Conn

func openNATS(addr string) (*nats.Conn, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	var opts []nats.Option
	if u.User != nil {
		opts = append(opts, nats.Secure(&tls.Config{InsecureSkipVerify: true}))
	}

	c, err := nats.Connect(addr, opts...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

var cMINIO *minio.Client

func openMINIO(addr string) (*minio.Client, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	var aKey, sKey string
	if u.User != nil {
		aKey = u.User.Username()
		sKey, _ = u.User.Password()
	}

	c, err := minio.New(u.Host, aKey, sKey, u.Scheme == "https")
	if err != nil {
		return nil, err
	}

	return c, nil
}
