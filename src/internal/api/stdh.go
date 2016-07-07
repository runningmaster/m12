package api

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"

	"internal/ctxutil"
	"internal/pref"
)

func stdh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if pref.Debug {
		if h, p := http.DefaultServeMux.Handler(r); p != "" {
			b := &stdhResponseWriter{rw: w}
			h.ServeHTTP(b, r)
			ctx = ctxutil.WithSize(ctx, int64(b.n))
		}
	} else {
		ctx = ctxutil.WithFail(ctx, fmt.Errorf("api: flag debug not found"))
	}
	*r = *r.WithContext(ctx)
}

type stdhResponseWriter struct {
	n  uint64
	rw http.ResponseWriter
}

func (w *stdhResponseWriter) Write(b []byte) (int, error) {
	n, err := w.rw.Write(b)
	atomic.AddUint64(&w.n, uint64(n))
	return n, err
}

func (w *stdhResponseWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *stdhResponseWriter) WriteHeader(statusCode int) {
	w.rw.WriteHeader(statusCode)
}

func (w *stdhResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.rw.(http.Hijacker).Hijack()
}

func (w *stdhResponseWriter) Count() uint64 {
	return atomic.LoadUint64(&w.n)
}
