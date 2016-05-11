package api

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

func pipeHead(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = withTime(ctx, time.Now())
		ctx = withUUID(ctx, nextUUID())
		ctx = withAddr(ctx, mineIP(r))
		h(ctx, w, r)
	}
}

func nextUUID() string {
	return fmt.Sprintf("%x", genUUID.Next())[:16]
}

func mineIP(r *http.Request) string {
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
