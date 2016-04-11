package api

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"internal/context/ctxutil"

	"github.com/rogpeppe/fastuuid"
	"golang.org/x/net/context"
)

var uuidPool = sync.Pool{
	New: func() interface{} {
		g, err := fastuuid.NewGenerator()
		if err != nil {
			return err
		}
		return g
	},
}

// GetGenerator gets generator from pool.
func getGenerator() (*fastuuid.Generator, error) {
	switch g := uuidPool.Get().(type) {
	case *fastuuid.Generator:
		return g, nil
	case error:
		return nil, g
	}

	return nil, fmt.Errorf("uuid: unreachable")
}

// PutGenerator puts generator back to the pool.
func putGenerator(x *fastuuid.Generator) {
	uuidPool.Put(x)
}

func pipeHead(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = ctxutil.WithTime(ctx, time.Now())
		ctx = ctxutil.WithIP(ctx, findIP(r))

		g, err := getGenerator()
		if err != nil {
			// FIXME log err
		}
		defer putGenerator(g)

		h(ctxutil.WithID(ctx, fmt.Sprintf("%x", g.Next())), w, r)
	}
}

func findIP(r *http.Request) string {
	var ip string
	if ip = r.Header.Get("X-Real-IP"); ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	ip, _, _ = net.SplitHostPort(ip)
	return ip
}
