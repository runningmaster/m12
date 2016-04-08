package core

import (
	"bytes"
	"fmt"
	"internal/net/s3"
	"net/http"
	"strings"

	"github.com/garyburd/redigo/redis"
	"golang.org/x/net/context"
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
type Handler func(context.Context, []byte) (interface{}, error)

// GetAuth gets auth(s).
func GetAuth(ctx context.Context, b []byte) (interface{}, error) {
	return applyOp(opGet, toAuth)(ctx, b)
}

// SetAuth sets auth(s).
func SetAuth(ctx context.Context, b []byte) (interface{}, error) {
	return applyOp(opSet, toAuth)(ctx, b)
}

// DelAuth deletes auth(s).
func DelAuth(ctx context.Context, b []byte) (interface{}, error) {
	return applyOp(opDel, toAuth)(ctx, b)
}

// GetLinkAddr gets linkAddr(s).
func GetLinkAddr(ctx context.Context, b []byte) (interface{}, error) {
	return applyOp(opGet, toLinkAddr)(ctx, b)
}

// SetLinkAddr sets linkAddr(s).
func SetLinkAddr(ctx context.Context, b []byte) (interface{}, error) {
	return applyOp(opSet, toLinkAddr)(ctx, b)
}

// DelLinkAddr deletes linkAddr(s).
func DelLinkAddr(ctx context.Context, b []byte) (interface{}, error) {
	return applyOp(opDel, toLinkAddr)(ctx, b)
}

// GetLinkDrug gets linkDrug(s).
func GetLinkDrug(ctx context.Context, b []byte) (interface{}, error) {
	return applyOp(opGet, toLinkDrug)(ctx, b)
}

// SetLinkDrug sets linkDrug(s).
func SetLinkDrug(ctx context.Context, b []byte) (interface{}, error) {
	return applyOp(opSet, toLinkDrug)(ctx, b)
}

// DelLinkDrug deletes linkDrug(s).
func DelLinkDrug(ctx context.Context, b []byte) (interface{}, error) {
	return applyOp(opDel, toLinkDrug)(ctx, b)
}

// GetLinkStat gets linkStat(s).
func GetLinkStat(ctx context.Context, b []byte) (interface{}, error) {
	return applyOp(opGet, toLinkStat)(ctx, b)
}

// SetLinkStat sets linkStat(s).
func SetLinkStat(ctx context.Context, b []byte) (interface{}, error) {
	return applyOp(opSet, toLinkStat)(ctx, b)
}

// DelLinkStat deletes linkStat(s).
func DelLinkStat(ctx context.Context, b []byte) (interface{}, error) {
	return applyOp(opDel, toLinkStat)(ctx, b)
}

func applyOp(op opType, to toType) Handler {
	return func(_ context.Context, b []byte) (interface{}, error) {
		var err error
		if b, err = mendGzipAndUTF8(b); err != nil {
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
func ToS3(ctx context.Context, b []byte) (interface{}, error) {
	if !strings.Contains(http.DetectContentType(b), "gzip") {
		return nil, fmt.Errorf("core: s3: gzip not found")
	}

	//err := s3.MkB("test")
	//if err != nil {
	//	return nil, err
	//}

	err := s3.Put("test", "name2", bytes.NewBuffer(b), "{}")
	if err != nil {
		return nil, err
	}

	return "OK", nil
}

// Ping calls Redis PING
func Ping(_ context.Context, b []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)

	return c.Do("PING")
}

// Info calls Redis INFO
func Info(_ context.Context, b []byte) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)

	b, err := redis.Bytes(c.Do("INFO"))
	if err != nil {
		return nil, err
	}

	return parseInfo(b)
}
