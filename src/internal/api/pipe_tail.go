package api

import (
	"net"
	"net/http"
	"time"

	"internal/context/ctxutil"
	"internal/log"

	"github.com/pivotal-golang/bytefmt"
	"golang.org/x/net/context"
)

func pipeTail(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		inf := informer{ctx, w, r}
		val := []interface{}{
			markEmpty(inf.id()[:16]),
			markEmpty(inf.address()),
			markEmpty(inf.method()),
			markEmpty(inf.path()),
			markEmpty(inf.auth()),
			inf.code(),
			bytefmt.ByteSize(uint64(inf.size())),
			markEmpty(inf.fail()),
			markEmpty(inf.time()),
			markEmpty(inf.agent()),
		}
		log.Println(val...)
		//if h != nil {
		//	h(ctx, w, r)
		//}
	}
}

func markEmpty(s string) string {
	if s != "" {
		return s
	}
	return "-"
}

type informer struct {
	c context.Context
	w http.ResponseWriter
	r *http.Request
}

func (i informer) path() string {
	var p string
	if p = i.r.URL.Path; p == "" {
		p = "/"
	}
	return p
}

func (i informer) method() string {
	return i.r.Method
}

func (i informer) address() string {
	var ip string
	if ip = i.r.Header.Get("X-Real-IP"); ip == "" {
		ip = i.r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = i.r.RemoteAddr
	}
	ip, _, _ = net.SplitHostPort(ip)
	return ip
}

func (i informer) agent() string {
	return i.r.UserAgent()
}

func (i informer) id() string {
	return ctxutil.IDFromContext(i.c)
}

func (i informer) auth() string {
	return ctxutil.AuthFromContext(i.c)
}

func (i informer) code() int {
	return ctxutil.CodeFromContext(i.c)
}

func (i informer) size() int64 {
	return ctxutil.SizeFromContext(i.c)
}

func (i informer) fail() string {
	if err := ctxutil.FailFromContext(i.c); err != nil {
		return err.Error()
	}
	return ""
}

func (i informer) time() string {
	if t := ctxutil.TimeFromContext(i.c); !t.IsZero() {
		return time.Since(t).String()
	}
	return ""
}
