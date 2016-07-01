package core

import (
	"internal/gzpool"

	"github.com/garyburd/redigo/redis"
)

const (
	keyZlog = "zset:meta"
)

func zlog(m jsonMeta) error {
	c := redisConn()
	defer closeConn(c)

	z, err := gzpool.Gzip(m.marshal())
	if err != nil {
		return err
	}

	_, err = c.Do("ZADD", keyZlog, m.Unix, z)
	return err
}

// GetZlog returns
func GetZlog(data []byte) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

	return redis.ByteSlices(c.Do("ZRANGEBYSCORE", keyZlog, "-inf", "+inf"))
	//if err != nil {
	//	return nil, err
	//}

	//out := make([]jsonMeta, len(res))
	//for i :=range res
}
