package core

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)

// RunC is "print-like" operation.
func RunC(cmd, base string) Handler {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
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

		var gsd redisGetSetDelOper
		if gsd, err = makeGetSetDeler(base, b); err != nil {
			return nil, err
		}

		return execGetSetDeler(cmd, gsd)
	}
}

func makeGetSetDeler(base string, b []byte) (redisGetSetDelOper, error) {
	switch base {
	case "auth":
		return decodeAuth(b), nil
	case "addr":
		return decodeLinkAddr(b), nil
	case "drug":
		return decodeLinkDrug(b), nil
	case "stat":
		return decodeLinkStat(b), nil
	}

	return nil, fmt.Errorf("core: unknown base %s", base)
}

func execGetSetDeler(cmd string, gsd redisGetSetDelOper) (interface{}, error) {
	c := redisGet()
	defer redisPut(c)

	switch cmd {
	case "get":
		return gsd.get(c)
	case "set":
		return gsd.set(c)
	case "del":
		return gsd.del(c)
	}

	return nil, fmt.Errorf("core: unknown command %s", cmd)
}
