package redis

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	cli    *nats.Conn
	logger = log.New(ioutil.Discard, "", log.LstdFlags)
)

func Run(addr string, log *log.Logger) error {
	if log != nil {
		logger = log
	}

	cli = &redis.Pool{
		Dial:        dial(addr),
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
	}

	var c redis.Conn
	if c = pool.Get(); c.Err() != nil {
		return fmt.Errorf("redis: %s", c.Err())
	}

	_ = c.Close()
	return nil
}

func dial(addr string) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		u, err := url.Parse(addr)
		if err != nil {
			return nil, err
		}

		c, err := redis.Dial("tcp", u.Host)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err != nil && c != nil {
				_ = c.Close()
			}
		}()

		if u.User != nil {
			if pw, ok := u.User.Password(); ok {
				_, err := c.Do("AUTH", pw)
				if err != nil {
					return nil, err
				}
			}
		}

		return c, nil
	}
}

func Get() redis.Conn {
	return cli.Get()
}

func Put(c io.Closer) error {
	return c.Close()
}

func ToInt64Safely(v interface{}) int64 {
	res, _ := redis.Int64(v, nil)
	return res
}

func ToStringSafely(v interface{}) string {
	res, _ := redis.String(v, nil)
	return res
}
