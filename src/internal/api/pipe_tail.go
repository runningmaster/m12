package api

import (
	"log"
	"net/http"
	"time"

	"github.com/pivotal-golang/bytefmt"
	"golang.org/x/net/context"
)

func pipeTail(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		inf := informer{ctx, w, r}
		val := []interface{}{
			markEmpty(inf.uuid()),
			markEmpty(inf.addr()),
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

func (i informer) addr() string {
	return addrFromCtx(i.c)
}

func (i informer) agent() string {
	return i.r.UserAgent()
}

func (i informer) uuid() string {
	return uuidFromCtx(i.c)
}

func (i informer) auth() string {
	return authFromCtx(i.c)
}

func (i informer) code() int64 {
	return codeFromCtx(i.c)
}

func (i informer) size() int64 {
	return sizeFromCtx(i.c)
}

func (i informer) fail() string {
	err := failFromCtx(i.c)
	if err != nil {
		return err.Error()
	}
	return ""
}

func (i informer) time() string {
	if t := timeFromCtx(i.c); !t.IsZero() {
		return time.Since(t).String()
	}
	return ""
}
