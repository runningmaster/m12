package redis

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	cli    *redis.Pool
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
	if c = cli.Get(); c.Err() != nil {
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

		if u.User != nil {
			if p, ok := u.User.Password(); ok {
				_, err = c.Do("AUTH", p)
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

func Put(c io.Closer) {
	if c != nil {
		_ = c.Close() // log ?
	}
}

func ToInt64Safely(v interface{}) int64 {
	res, _ := redis.Int64(v, nil)
	return res
}

func ToStringSafely(v interface{}) string {
	res, _ := redis.String(v, nil)
	return res
}

// TODO: parse keyspace
//	"keyspace": {
//		"db0": "keys=1,expires=0,avg_ttl=0"
//	},
func InfoToJSON(reply interface{}, err error) (interface{}, error) {
	b, err := redis.Bytes(reply, err)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	mapper := make(map[string]map[string]string)

	var (
		line  string
		sect  string
		split []string
	)

	for scanner.Scan() {
		line = strings.ToLower(scanner.Text())
		if line == "" {
			continue
		}
		if strings.Contains(line, "#") {
			sect = line[2:]
			mapper[sect] = make(map[string]string)
			continue
		}
		split = strings.Split(line, ":")
		mapper[sect][split[0]] = split[1]
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return mapper, nil
}
