package core

import (
	"bytes"
	"fmt"
	"net/http"

	"internal/net/s3"

	"github.com/garyburd/redigo/redis"
)

// Handler is func for processing data from api.
type Handler func(r *http.Request) (interface{}, error)

// RunC is "print-like" operation.
func RunC(cmd, base string) Handler {
	return func(r *http.Request) (interface{}, error) {
		var (
			b   []byte
			err error
		)
		if b, err = readBody(r); err != nil {
			return nil, err
		}

		if b, err = mendIfGzip(b); err != nil {
			return nil, err
		}

		if b, err = mendIfUTF8(b); err != nil {
			return nil, err
		}

		var gsd getsetdeler
		if gsd, err = makeGetSetDeler(base, b); err != nil {
			return nil, err
		}

		return execGetSetDeler(cmd, gsd)
	}
}

func makeGetSetDeler(base string, b []byte) (getsetdeler, error) {
	switch base {
	case "auth":
		return decodeAuth(b), nil
	case "addr":
		return decodeLinkAddr(b), nil
	case "drug":
		return decodeLinkDrug(b), nil
	case "stat":
		return decodeLinkStat(b), nil
	}

	return nil, fmt.Errorf("core: unknown base %s", base)
}

func execGetSetDeler(cmd string, gsd getsetdeler) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)

	switch cmd {
	case "get":
		return gsd.get(c)
	case "set":
		return gsd.set(c)
	case "del":
		return gsd.del(c)
	}

	return nil, fmt.Errorf("core: unknown command %s", cmd)
}

// ToS3 sends data to s3 interface
func ToS3(r *http.Request) (interface{}, error) {
	var (
		b   []byte
		err error
	)
	if b, err = readBody(r); err != nil {
		return nil, err
	}

	if !isTypeGzip(b) {
		return nil, fmt.Errorf("core: s3: gzip not found")
	}

	//err := s3.MkB("test")
	//if err != nil {
	//	return nil, err
	//}

	if err = s3.Put("test", "name2", bytes.NewBuffer(b), "{}"); err != nil {
		return nil, err
	}

	return "OK", nil
}

// Ping calls Redis PING
func Ping(_ *http.Request) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)

	return c.Do("PING")
}

// Info calls Redis INFO
func Info(_ *http.Request) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)

	b, err := redis.Bytes(c.Do("INFO"))
	if err != nil {
		return nil, err
	}

	return parseInfo(b)
}
