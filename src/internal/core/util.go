package core

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"internal/compress/gzutil"

	"github.com/garyburd/redigo/redis"
	"github.com/spkg/bom"
)

type redisGetSetDelOper interface {
	get(redis.Conn) ([]interface{}, error)
	set(redis.Conn) (interface{}, error)
	del(redis.Conn) (interface{}, error)
}

func toInt64(v interface{}) int64 {
	res, _ := redis.Int64(v, nil)
	return res
}

func toString(v interface{}) string {
	res, _ := redis.String(v, nil)
	return res
}

func stringOK() string {
	return http.StatusText(http.StatusOK)
}

func readBody(r *http.Request) ([]byte, error) {
	defer func(c io.Closer) {
		_ = c.Close()
	}(r.Body)

	return ioutil.ReadAll(r.Body)
}

func isTypeGzip(b []byte) bool {
	return gzutil.IsGzipInString(http.DetectContentType(b))
}

func isTypeUTF8(b []byte) bool {
	return strings.Contains(http.DetectContentType(b), "text/plain; charset=utf-8")
}

func mendIfGzip(b []byte) ([]byte, error) {
	if isTypeGzip(b) {
		unz, err := gzutil.Gunzip(b)
		if err != nil {
			return nil, err
		}
		return unz, nil
	}

	return b, nil
}

func mendIfUTF8(b []byte) ([]byte, error) {
	if isTypeUTF8(b) {
		return bom.Clean(b), nil
	}

	return b, nil
}
