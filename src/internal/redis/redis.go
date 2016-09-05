package redis

import (
	"log"
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Init return active pool connections to REDIS Server
func Init(addr string) (*redis.Pool, error) {
	return makePool(addr)
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
