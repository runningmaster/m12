package core

import (
	"internal/errors"

	"github.com/garyburd/redigo/redis"
)

// Ping calls Redis PING
func Ping(b []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)

	return c.Do("PING")
}

// Info calls Redis INFO
func Info(b []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)

	b, err := redis.Bytes(c.Do("INFO"))
	if err != nil {
		return nil, errors.Locus(err)
	}

	return parseInfo(b)
}

func GetAuth(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeAuth(data).get(c)
}

func SetAuth(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeAuth(data).set(c)
}

func DelAuth(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeAuth(data).del(c)
}

func GetLinkAddr(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkAddr(data).get(c)
}

func SetLinkAddr(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkAddr(data).set(c)
}

func DelLinkAddr(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkAddr(data).del(c)
}

func GetLinkDrug(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkDrug(data).get(c)
}

func SetLinkDrug(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkDrug(data).set(c)
}

func DelLinkDrug(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkDrug(data).del(c)
}

func GetLinkStat(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkStat(data).get(c)
}

func SetLinkStat(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkStat(data).set(c)
}

func DelLinkStat(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkStat(data).del(c)
}
