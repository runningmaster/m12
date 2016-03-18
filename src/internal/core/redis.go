package core

import (
	"io"
	"net/url"
	"time"

	"internal/errors"
	"internal/flag"
	"internal/server"

	"github.com/garyburd/redigo/redis"
)

var redisPool rediser

type (
	// Rediser is dead simple common interface for executing commands to redis
	rediser interface {
		Get() redis.Conn
		Put(c io.Closer)

		//	ExecOne([]interface{}) (interface{}, error)
		//	ExecMulti([][]interface{}) ([]interface{}, error)
	}

	redisServer struct {
		pool interface {
			Get() redis.Conn
		}
	}
)

func init() {
	var err error
	redisPool, err = newRedis(flag.Redis)
	if err != nil {
		server.FailFast = errors.Locus(err)
	}
}

// New creates a server configured from default values.
func newRedis(addr string) (rediser, error) {
	pool := &redis.Pool{
		Dial:        dial(addr),
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
	}

	var c redis.Conn
	if c = pool.Get(); c.Err() != nil {
		return nil, errors.Locus(c.Err())
	}
	_ = c.Close()

	return &redisServer{pool: pool}, nil
}

func dial(addr string) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		u, err := url.Parse(addr)
		if err != nil {
			return nil, errors.Locus(err)
		}

		c, err := redis.Dial("tcp", u.Host)
		if err != nil {
			return nil, errors.Locus(err)
		}

		defer func() {
			if err != nil && c != nil {
				_ = c.Close()
			}
		}()

		if u.User != nil {
			if pw, ok := u.User.Password(); ok {
				if _, err := c.Do("AUTH", pw); err != nil {
					return nil, errors.Locus(err)
				}
			}
		}

		return c, nil
	}
}

func (s *redisServer) Get() redis.Conn {
	return s.pool.Get()
}

func (s *redisServer) Put(c io.Closer) {
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
