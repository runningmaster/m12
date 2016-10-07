package pipe

import (
	"fmt"
	"net"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/rogpeppe/fastuuid"
)

var genUUID = fastuuid.MustNewGenerator()

func Head(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = withUUID(ctx, nextUUID())
		ctx = withHost(ctx, mineHost(r))
		ctx = withUser(ctx, mineUser(r))
		ctx = withTime(ctx, time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
