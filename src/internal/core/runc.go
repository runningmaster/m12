package core

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)

// RunC is "print-like" operation.
func RunC(cmd, base string) Handler {
	return func(_ context.Context, _ http.ResponseWriter, r *http.Request) (interface{}, error) {
		b, err := readClose(r.Body)
		if err != nil {
			return nil, err
		}

		b, err = mendIfGzip(b)
		if err != nil {
			return nil, err
		}

		b, err = mendIfUTF8(b)
		if err != nil {
			return nil, err
		}

		v, err := makeGetSetDeler(base, b)
		if err != nil {
			return nil, err
		}

		return execGetSetDeler(cmd, v)
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
