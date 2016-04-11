package redispool

import (
	"io"
	"net/url"
	"time"

	"internal/flag"
	"internal/log"

	"github.com/garyburd/redigo/redis"
)

var redisServer connGetter

type (
	connGetter interface {
		Get() redis.Conn
	}
)

func init() {
	var err error
	redisServer, err = newRedis(flag.Redis)
	if err != nil {
		log.Fatal(err)
	}
}

// New creates a server configured from default values.
func newRedis(addr string) (connGetter, error) {
	pool := &redis.Pool{
		Dial:        dial(addr),
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
	}

	var c redis.Conn
	if c = pool.Get(); c.Err() != nil {
		return nil, c.Err()
	}
	_ = c.Close()

	return pool, nil
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
				if _, err := c.Do("AUTH", pw); err != nil {
					return nil, err
				}
			}
		}

		return c, nil
	}
}

//
func Get() redis.Conn {
	return redisServer.Get()
}

//
func Put(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}
