package core

import (
	"crypto/tls"
	"io"
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
	minio "github.com/minio/minio-go"
	"github.com/nats-io/nats"
)

var cNATS *nats.Conn

func openNATS(addr string) (*nats.Conn, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	opts := []nats.Option{nats.MaxReconnects(-1)}
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

var pREDIS *redis.Pool

func openREDIS(addr string) (*redis.Pool, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	p := &redis.Pool{
		MaxIdle:     128,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", u.Host)
		},
	}

	c := p.Get()
	defer closeConn(c)

	return p, c.Err()
}

func redisConn() redis.Conn {
	return pREDIS.Get()
}

func closeConn(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}
