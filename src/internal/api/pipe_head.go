package api

import (
	"net"
	"net/http"
	"time"

	"internal/context/ctxutil"
	"internal/crypto/uuid"

	"golang.org/x/net/context"
)

func pipeHead(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = ctxutil.WithTime(ctx, time.Now())
		ctx = ctxutil.WithID(ctx, uuid.Next())
		ctx = ctxutil.WithIP(ctx, findIP(r))
		h(ctx, w, r)
	}
}

func findIP(r *http.Request) string {
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
