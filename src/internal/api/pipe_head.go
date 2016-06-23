package api

import (
	"fmt"
	"net"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/rogpeppe/fastuuid"
)

var genUUID = fastuuid.MustNewGenerator()

func pipeHead(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx = ctxWithUUID(ctx, nextUUID())
		ctx = ctxWithHost(ctx, mineHost(r))
		ctx = ctxWithUser(ctx, mineUser(r))
		ctx = ctxWithTime(ctx, time.Now())
		h(ctx, w, r)
	}
}

func nextUUID() string {
	return fmt.Sprintf("%x", genUUID.Next())
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

func mineUser(r *http.Request) string {
	u := r.UserAgent()
	if !utf8.Valid([]byte(u)) {
		u = fmt.Sprintf("[Warning: non UTF-8]: %s", u)
	}
	return u
}
