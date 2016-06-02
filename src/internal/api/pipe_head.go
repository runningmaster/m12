package api

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rogpeppe/fastuuid"
	"golang.org/x/net/context"
)

var genUUID = fastuuid.MustNewGenerator()

func pipeHead(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = ctxWithUUID(ctx, nextUUID())
		ctx = ctxWithHost(ctx, mineHost(r))
		ctx = ctxWithTime(ctx, time.Now())
		h(ctx, w, r)
	}
}

func nextUUID() string {
	return fmt.Sprintf("%x", genUUID.Next())[:16]
}

func mineHost(r *http.Request) string {
	h := r.Header.Get("X-Forwarded-For")
	if h != "" {
		return h
	}

	h = r.Header.Get("X-Real-IP")
	if h != "" {
		return h
	}

	var err error
	h, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		h = "?.?.?.?"
	}

	return h
}
