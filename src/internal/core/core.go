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

// GetAuth returns slice of things or nil instead ones (if thing doesn't exist).
func GetAuth(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeAuth(data).get(c)
}

// SetAuth upserts things and returns simple "OK".
func SetAuth(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeAuth(data).set(c)
}

// DelAuth deletes things and returns the number of things that were removed.
func DelAuth(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeAuth(data).del(c)
}

// GetLinkAddr returns slice of things or nil instead ones (if thing doesn't exist).
func GetLinkAddr(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkAddr(data).get(c)
}

// SetLinkAddr upserts things and returns simple "OK".
func SetLinkAddr(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkAddr(data).set(c)
}

// DelLinkAddr deletes things and returns the number of things that were removed.
func DelLinkAddr(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkAddr(data).del(c)
}

// GetLinkDrug returns slice of things or nil instead ones (if thing doesn't exist).
func GetLinkDrug(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkDrug(data).get(c)
}

// SetLinkDrug upserts things and returns simple "OK".
func SetLinkDrug(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkDrug(data).set(c)
}

// DelLinkDrug deletes things and returns the number of things that were removed.
func DelLinkDrug(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkDrug(data).del(c)
}

// GetLinkStat returns slice of things or nil instead ones (if thing doesn't exist).
func GetLinkStat(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkStat(data).get(c)
}

// SetLinkStat upserts things and returns simple "OK".
func SetLinkStat(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkStat(data).set(c)
}

// DelLinkStat deletes things and returns the number of things that were removed.
func DelLinkStat(data []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)
	return decodeLinkStat(data).del(c)
}
