package core

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"internal/flag"

	"github.com/garyburd/redigo/redis"
)

var redisServer connGetter

type (
	connGetter interface {
		Get() redis.Conn
	}

	redisGetSetDelOper interface {
		get(redis.Conn) ([]interface{}, error)
		set(redis.Conn) (interface{}, error)
		del(redis.Conn) (interface{}, error)
	}
)

func initRedis() error {
	var err error
	redisServer, err = newRedis(flag.Redis)
	if err != nil {
		return fmt.Errorf("core: redis: %s", err)
	}

	return nil
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
				_, err := c.Do("AUTH", pw)
				if err != nil {
					return nil, err
				}
			}
		}

		return c, nil
	}
}

func redisGet() redis.Conn {
	return redisServer.Get()
}

func redisPut(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}

func toInt64(v interface{}) int64 {
	res, _ := redis.Int64(v, nil)
	return res
}

func toString(v interface{}) string {
	res, _ := redis.String(v, nil)
	return res
}
