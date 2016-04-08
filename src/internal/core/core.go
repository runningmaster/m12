package core

import (
	"bytes"
	"fmt"
	"net/http"

	"internal/net/s3"

	"github.com/garyburd/redigo/redis"
)

type (
	toType int
	opType int
)

const (
	toAuth toType = iota
	toLinkAddr
	toLinkDrug
	toLinkStat

	opGet opType = iota
	opSet
	opDel
)

// Handler is func for processing data from api.
type Handler func(r *http.Request) (interface{}, error)

// GetAuth gets auth(s).
func GetAuth(r *http.Request) (interface{}, error) {
	return applyOp(opGet, toAuth)(r)
}

// SetAuth sets auth(s).
func SetAuth(r *http.Request) (interface{}, error) {
	return applyOp(opSet, toAuth)(r)
}

// DelAuth deletes auth(s).
func DelAuth(r *http.Request) (interface{}, error) {
	return applyOp(opDel, toAuth)(r)
}

// GetLinkAddr gets linkAddr(s).
func GetLinkAddr(r *http.Request) (interface{}, error) {
	return applyOp(opGet, toLinkAddr)(r)
}

// SetLinkAddr sets linkAddr(s).
func SetLinkAddr(r *http.Request) (interface{}, error) {
	return applyOp(opSet, toLinkAddr)(r)
}

// DelLinkAddr deletes linkAddr(s).
func DelLinkAddr(r *http.Request) (interface{}, error) {
	return applyOp(opDel, toLinkAddr)(r)
}

// GetLinkDrug gets linkDrug(s).
func GetLinkDrug(r *http.Request) (interface{}, error) {
	return applyOp(opGet, toLinkDrug)(r)
}

// SetLinkDrug sets linkDrug(s).
func SetLinkDrug(r *http.Request) (interface{}, error) {
	return applyOp(opSet, toLinkDrug)(r)
}

// DelLinkDrug deletes linkDrug(s).
func DelLinkDrug(r *http.Request) (interface{}, error) {
	return applyOp(opDel, toLinkDrug)(r)
}

// GetLinkStat gets linkStat(s).
func GetLinkStat(r *http.Request) (interface{}, error) {
	return applyOp(opGet, toLinkStat)(r)
}

// SetLinkStat sets linkStat(s).
func SetLinkStat(r *http.Request) (interface{}, error) {
	return applyOp(opSet, toLinkStat)(r)
}

// DelLinkStat deletes linkStat(s).
func DelLinkStat(r *http.Request) (interface{}, error) {
	return applyOp(opDel, toLinkStat)(r)
}

func applyOp(op opType, to toType) Handler {
	return func(r *http.Request) (interface{}, error) {
		var (
			b   []byte
			err error
		)
		if b, err = readBody(r); err != nil {
			return nil, err
		}

		if b, err = mendIfGzip(b); err != nil {
			return nil, err
		}

		if b, err = mendIfUTF8(b); err != nil {
			return nil, err
		}

		var gsd getsetdeler
		if gsd, err = makeGetSetDeler(to, b); err != nil {
			return nil, err
		}

		return execGetSetDeler(gsd, op)
	}
}

func makeGetSetDeler(to toType, b []byte) (getsetdeler, error) {
	switch to {
	case toAuth:
		return decodeAuth(b), nil
	case toLinkAddr:
		return decodeLinkAddr(b), nil
	case toLinkDrug:
		return decodeLinkDrug(b), nil
	case toLinkStat:
		return decodeLinkStat(b), nil
	}

	return nil, fmt.Errorf("core: unknown sys type %v", to)
}

func execGetSetDeler(gsd getsetdeler, op opType) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)

	switch op {
	case opGet:
		return gsd.get(c)
	case opSet:
		return gsd.set(c)
	case opDel:
		return gsd.del(c)
	}

	return nil, fmt.Errorf("core: unknown op type %v", op)
}

// ToS3 sends data to s3 interface
func ToS3(r *http.Request) (interface{}, error) {
	var (
		b   []byte
		err error
	)
	if b, err = readBody(r); err != nil {
		return nil, err
	}

	if !isTypeGzip(b) {
		return nil, fmt.Errorf("core: s3: gzip not found")
	}

	//err := s3.MkB("test")
	//if err != nil {
	//	return nil, err
	//}

	if err = s3.Put("test", "name2", bytes.NewBuffer(b), "{}"); err != nil {
		return nil, err
	}

	return "OK", nil
}

// Ping calls Redis PING
func Ping(_ *http.Request) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)

	return c.Do("PING")
}

// Info calls Redis INFO
func Info(_ *http.Request) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)

	b, err := redis.Bytes(c.Do("INFO"))
	if err != nil {
		return nil, err
	}

	return parseInfo(b)
}
