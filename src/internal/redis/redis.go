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

const statusOK = "OK"

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

func getConn() redis.Conn {
	return cli.Get()
}

func putConn(c io.Closer) {
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

func Ping() (interface{}, error) {
	c := getConn()
	defer putConn(c)

	return c.Do("PING")
}

// TODO: parse keyspace
//	"keyspace": {
//		"db0": "keys=1,expires=0,avg_ttl=0"
//	},
func Info() (interface{}, error) {
	c := getConn()
	defer putConn(c)

	b, err := redis.Bytes(c.Do("INFO"))
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

func ConvFromInt64s(src ...int64) []interface{} {
	dst := make([]interface{}, len(src))
	for i := range dst {
		dst[i] = src[i]
	}
	return dst
}

func ConvFromInt64sWithKey(key string, src ...int64) []interface{} {
	dst := make([]interface{}, len(src)+1)
	dst[0] = key
	for i, j := 0, 1; i < len(dst)-1; i, j = i+1, j+1 {
		dst[j] = src[i]
	}
	return dst
}

func ConvFromStrings(src ...string) []interface{} {
	dst := make([]interface{}, len(src))
	for i := range dst {
		dst[i] = src[i]
	}
	return dst
}

func ConvFromStringsWithKey(key string, src ...string) []interface{} {
	dst := make([]interface{}, len(src)+1)
	dst[0] = key
	for i, j := 0, 1; i < len(dst)-1; i, j = i+1, j+1 {
		dst[j] = src[i]
	}
	return dst
}

// HMSET is wrapper func and returns "simple string reply". Key must be first in array.
func HMSET(keyAndFieldVals ...interface{}) (interface{}, error) {
	c := getConn()
	defer putConn(c)

	return c.Do("HMSET", keyAndFieldVals...)
}

// HMGET is wrapper func and returns "array reply". Key must be first in array.
func HMGET(keyAndFields ...interface{}) ([]interface{}, error) {
	c := getConn()
	defer putConn(c)

	return redis.Values(c.Do("HMGET", keyAndFields...))
}

// HDEL is wrapper func and returns "integer reply". Key must be first in array.
func HDEL(keyAndFields ...interface{}) (interface{}, error) {
	c := getConn()
	defer putConn(c)

	return c.Do("HDEL", keyAndFields...)
}

// HMSETM is wrapper func and returns "simple string reply". Key must be first in array.
func HMSETM(keyAndFieldVals ...[]interface{}) (interface{}, error) {
	c := getConn()
	defer putConn(c)

	var err error
	for i := range keyAndFieldVals {
		if len(keyAndFieldVals[i]) == 0 {
			continue
		}

		err = c.Send("DEL", keyAndFieldVals[i][0])
		if err != nil {
			return nil, err
		}

		err = c.Send("HMSET", keyAndFieldVals[i]...)
		if err != nil {
			return nil, err
		}
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	return statusOK, nil
}

// HMGETM is wrapper func and returns "array reply". Key must be first in array.
func HMGETM(keyAndFields ...[]interface{}) ([][]interface{}, error) {
	c := getConn()
	defer putConn(c)

	var err error
	for i := range keyAndFields {
		err = c.Send("HMGET", keyAndFields[i]...)
		if err != nil {
			return nil, err
		}
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	out := make([][]interface{}, len(keyAndFields))
	var res interface{}
	for i := 0; i < len(out); i++ {
		res, err = c.Receive()
		if err != nil {
			return nil, err
		}
		out[i], err = redis.Values(res, err)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

// DEL is wrapper func and returns "integer reply".
func DEL(keys ...interface{}) (interface{}, error) {
	c := getConn()
	defer putConn(c)

	return c.Do("DEL", keys...)
}

// SMEMBERS is wrapper func and returns "array reply".
func SMEMBERS(key interface{}) ([]interface{}, error) {
	c := getConn()
	defer putConn(c)

	return redis.Values(c.Do("SMEMBERS", key))
}

// SISMEMBER is wrapper func and returns "integer reply". Key must be first in array.
func SISMEMBER(key, member interface{}) (interface{}, error) {
	c := getConn()
	defer putConn(c)

	return c.Do("SISMEMBER", key, member)
}

// SISMEMBERM is wrapper func and returns "array reply". Key must be first in array.
func SISMEMBERM(keyAndMembers ...interface{}) ([]interface{}, error) {
	if len(keyAndMembers) == 0 {
		return nil, fmt.Errorf("no arguments")
	}

	c := getConn()
	defer putConn(c)

	key := keyAndMembers[0]
	var err error
	for i := 1; i < len(keyAndMembers); i++ {
		err = c.Send("SISMEMBER", key, keyAndMembers[i])
		if err != nil {
			return nil, err
		}
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	out := make([]interface{}, len(keyAndMembers)-1)
	var res interface{}
	for i := 0; i < len(out); i++ {
		res, err = c.Receive()
		if err != nil {
			return nil, err
		}
		out[i] = res
	}

	return out, nil
}

// SADD is wrapper func and returns "integer reply". Key must be first in array.
func SADD(keyAndMembers ...interface{}) (interface{}, error) {
	c := getConn()
	defer putConn(c)

	return c.Do("SADD", keyAndMembers...)
}

// SREM is wrapper func and returns "integer reply". Key must be first in array.
func SREM(keyAndMembers ...interface{}) (interface{}, error) {
	c := getConn()
	defer putConn(c)

	return c.Do("SREM", keyAndMembers...)
}
