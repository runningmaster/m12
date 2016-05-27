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
		ctx = ctxWithTime(ctx, time.Now())
		ctx = ctxWithUUID(ctx, nextUUID())
		ctx = ctxWithAddr(ctx, mineAddr(r))
		h(ctx, w, r)
	}
}

func nextUUID() string {
	return fmt.Sprintf("%x", genUUID.Next())[:16]
}

func mineAddr(r *http.Request) string {
	var ip string
	if ip = r.Header.Get("X-Forwarded-For"); ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}

	if ip == "" {
		ip = r.RemoteAddr
	}
	ip, _, _ = net.SplitHostPort(ip)

	return ip
}
