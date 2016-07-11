package api

import (
	"encoding/json"

	"internal/gzip"

	"github.com/garyburd/redigo/redis"
)

const (
	keyZlog = "zset:meta"
)

func zlog(m jsonMeta) error {
	c := redisConn()
	defer closeConn(c)

	z, err := gzip.Compress(m.marshal())
	if err != nil {
		return err
	}

	err = c.Send("ZADD", keyZlog, m.Unix, m.UUID)
	if err != nil {
		return err
	}

	err = c.Send("SET", m.UUID, z, "EX", 60*60*24*3)
	if err != nil {
		return err
	}

	return c.Flush()
}

func getZlog(data []byte) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

	res, err := redis.Strings(c.Do("ZRANGEBYSCORE", keyZlog, "-inf", "+inf"))
	if err != nil {
		return nil, err
	}
	for i := range res {
		err = c.Send("GET", res[i])
		if err != nil {
			return nil, err
		}
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	out := make([]jsonMeta, 0, len(res))
	var r []byte
	var z []byte
	var m jsonMeta
	for range res {
		z, err = redis.Bytes(c.Receive())
		if err != nil && err != redis.ErrNil {
			return nil, err
		}

		r, err = gzip.Uncompress(z)
		if err != nil {
			return nil, err
		}

		m, err = unmarshalMeta(r)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}

	return out, nil
}

func getMeta(data []byte) (interface{}, error) {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	c := redisConn()
	defer closeConn(c)

	z, err := redis.Bytes(c.Do("GET", v))
	if err != nil /*&& err != redis.ErrNil*/ {
		return nil, err
	}

	r, err := gzip.Uncompress(z)
	if err != nil {
		return nil, err
	}

	m, err := unmarshalMeta(r)
	if err != nil {
		return nil, err
	}

	return m, nil
}
