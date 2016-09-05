package nats

import (
	"crypto/tls"
	"net/url"

	"github.com/nats-io/nats"
)

// Init return active connection to NATS Server
func Init(addr string) (*nats.Conn, error) {
	return makeConn(addr)
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

/*
	_, err = cNATS.Subscribe(subjectSteamIn, func(m *nats.Msg) {
		go proc(m.Data)
	})
	if err != nil {
		return err
	}

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
