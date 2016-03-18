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

// ExecMulti is wrapper for chain
// Send(commandName string, args ...interface{}) error
// Flush() error
// Receive() (reply interface{}, err error)
func (s *redisServer) ExecMulti(args [][]interface{}) ([]interface{}, error) {
	c := s.pool.Get()
	defer func(c io.Closer) {
		_ = c.Close()
	}(c)

	res := make([]interface{}, 0, len(args))
	var count int
	for i := range args {
		if cmd, ok := args[i][0].(string); ok {
			count++
			var err error
			switch len(args[i]) {
			case 1:
				err = c.Send(cmd)
			default:
				err = c.Send(cmd, args[i][1:]...)
			}
			if err != nil {
				return nil, errors.Locusf("database/redis: send() fails after %d attempts", count)
			}
		}
	}

	err := c.Flush()
	if err != nil {
		return nil, errors.Locusf("database/redis: flush() fails  after %d send() attempts", count)
	}

	for i := 0; i < count; i++ {
		rcv, err := c.Receive()
		if err != nil {
			return nil, errors.Locusf("database/redis: receive() fails after %d attempts", i)
		}
		res = append(res, rcv)
	}

	return res, nil
}

func toInt64(v interface{}) int64 {
	res, _ := redis.Int64(v, nil)
	return res
}

func toString(v interface{}) string {
	res, _ := redis.String(v, nil)
	return res
}
