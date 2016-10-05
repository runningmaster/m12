package nats

import (
	"crypto/tls"
	"net/url"

	"github.com/nats-io/nats"
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

/*


	go publishing(bucketStreamOut, subjectSteamOut, listN, tickD)
	go publishing(bucketStreamIn, subjectSteamIn, listN, tickD)

*/

/*
func publishing(bucket, subject string, n int, d time.Duration) {
	var err error
	for range time.Tick(d) {
		err = publish(bucket, subject, n)
		if err != nil {
			log.Println(err)
		}
	}
}

func publish(bucket, subject string, n int) error {
	l, err := listObjectsN(bucket, n)
	if err != nil {
		return err
	}

	m := make([][]byte, len(l))
	for i := range l {
		m[i] = pair{bucket, l[i]}.marshal()
	}

	return publishEach(subject, m...)
}

func publishEach(subject string, msgs ...[]byte) error {
	var err error
	for i := range msgs {
		err = cNATS.Publish(subject, msgs[i])
		if err != nil {
			return err
		}
	}
	return nil
}
*/
