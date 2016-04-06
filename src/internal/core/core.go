package core

import (
	"internal/errors"
	"internal/flag"

	"github.com/garyburd/redigo/redis"
	"golang.org/x/net/context"
)

type SysType int

const (
	SysAuth SysType = iota
	SysLinkAddr
	SysLinkDrug
	SysLinkStat
)

// Handler is func for processing data from api
type Handler func(context.Context, []byte) (interface{}, error)

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
		return nil, errors.Locus(err)
	}

	return parseInfo(b)
}

// SysOp is wrapper for get/sed/del ops
func SysOp(st SysType, op flag.OpType) Handler {
	return func(_ context.Context, b []byte) (interface{}, error) {
		c := redisPool.Get()
		defer redisPool.Put(c)

		var gsd getsetdeler
		switch st {
		case SysAuth:
			gsd = decodeAuth(b)
		case SysLinkAddr:
			gsd = decodeLinkAddr(b)
		case SysLinkDrug:
			gsd = decodeLinkDrug(b)
		case SysLinkStat:
			gsd = decodeLinkStat(b)
		default:
			return nil, errors.Locusf("core: unknown sys type %v", st)
		}

		switch op {
		case flag.OpGet:
			return gsd.get(c)
		case flag.OpSet:
			return gsd.set(c)
		case flag.OpDel:
			return gsd.del(c)
		}

		return nil, errors.Locusf("core: unknown op type %v", op)
	}
}

// ToS3 sends data to s3 interface
func ToS3(ctx context.Context, data []byte) (interface{}, error) {
	// TODO: fail fast if not gzip (!)
	return "ToS3 FIXME", nil
}
