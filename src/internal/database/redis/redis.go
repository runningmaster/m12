package redis

import (
	"io"
	"log"
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
	Conv Converter = conv{}
)

type Converter interface {
	ToInt64(interface{}, error) (int64, error)
	ToString(interface{}, error) (string, error)
	ToStrings(interface{}, error) ([]string, error)
	ToIntfs(interface{}, error) ([]interface{}, error)
	ToBytes(interface{}, error) ([]byte, error)
	NotErrNil(error) bool
}

// Init inits client for REDIS Server
func Init(addr string) error {
	p, err := makePool(addr)
	if err != nil {
		return err
	}
	pool = p
	return nil
}

func makePool(addr string) (*redis.Pool, error) {
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

	return p, waitDBFromDisk(p, 1*time.Second)
}

func waitDBFromDisk(p *redis.Pool, d time.Duration) error {
	c := p.Get()
	defer c.Close()

	t := time.NewTicker(d)
	var err error
	for range t.C {
		_, err = c.Do("PING")
		if err != nil {
			log.Println(err)
			continue
		}
		break
	}
	t.Stop()
	return err
}

func Conn() redis.Conn {
	return pool.Get()
}

func Free(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}

type conv struct{}

func (conv) ToInt64(v interface{}, err error) (int64, error) {
	return redis.Int64(v, err)
}

func (conv) ToString(v interface{}, err error) (string, error) {
	return redis.String(v, err)
}

func (conv) ToStrings(v interface{}, err error) ([]string, error) {
	return redis.Strings(v, err)
}

func (conv) ToIntfs(v interface{}, err error) ([]interface{}, error) {
	return redis.Values(v, err)
}

func (conv) ToBytes(v interface{}, err error) ([]byte, error) {
	return redis.Bytes(v, err)
}

func (conv) NotErrNil(err error) bool {
	return err != redis.ErrNil
}
