package core

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"internal/gzip"

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

func TestStreamOut() error {
	var err error

	cNATS, err = openNATS("nats://Morion:0258790@195.128.18.66:4222")
	if err != nil {
		return err
	}

	_, err = cNATS.Subscribe(subjectSteamOut, func(m *nats.Msg) {
		go procTest(m.Data)
	})
	if err != nil {
		return err
	}

	return nil
}

func procTest(pair []byte) {
	start := time.Now()
	req, err := http.NewRequest("POST", "https://195.128.18.66/stream/get-data", bytes.NewReader(pair))
	if err != nil {
		log.Print(fmt.Errorf("DEBUG 1: %v", err))
		return
	}
	req.SetBasicAuth("api", "key-sysdba")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := cli.Do(req)
	if err != nil {
		log.Print(fmt.Errorf("DEBUG 2: %v", err))
		return
	}
	defer res.Body.Close()

	var meta []byte
	meta, err = base64.StdEncoding.DecodeString(res.Header.Get("Content-Meta"))
	if err != nil {
		log.Print(fmt.Errorf("DEBUG 3: %v", err))
		return
	}

	data := new(bytes.Buffer)
	_, err = io.Copy(data, res.Body)
	if err != nil {
		log.Print(fmt.Errorf("DEBUG 4: %v", err))
		return
	}

	if res.StatusCode >= 300 {
		log.Print(fmt.Errorf("POST failed with code %d: %v", res.StatusCode, data.String()))
		return
	}

	_, object, err := unmarshaPairExt(pair)
	if err != nil {
		log.Print(fmt.Errorf("DEBUG 5: %v", err))
		return
	}

	s := strings.Split(object, "_")
	if len(s) == 0 {
		log.Print(fmt.Errorf("Invalid object %v", object))
		return
	}

	m, err := unmarshalMeta(meta)
	if err != nil {
		log.Print(fmt.Errorf("DEBUG 6: %v", err))
		return
	}

	d, err := gzip.Uncompress(data.Bytes())
	if err != nil {
		log.Print(fmt.Errorf("DEBUG 7: %v", err))
		return
	}

	var v []struct{}
	err = json.Unmarshal(d, &v)
	if err != nil {
		log.Print(fmt.Errorf("DEBUG 8: %v", err))
		return
	}

	log.Printf("%s %s %t: %d %s %s", s[0], m.UUID, strings.Contains(m.UUID, s[0]), len(v), m.Proc, time.Since(start).String())
}
