package api

import (
	"net"
	"net/http"
	"time"

	"internal/crypto/uuid"

	"golang.org/x/net/context"
)

func pipeHead(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = withTime(ctx, time.Now())
		ctx = withUUID(ctx, uuid.Next())
		ctx = withAddr(ctx, mineIP(r))
		h(ctx, w, r)
	}
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
