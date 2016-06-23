package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pivotal-golang/bytefmt"
)

const magicLen = 8

func pipeTail(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inf := informer{ctx, w, r}
		log.Println( // log.New() ?
			inf.code(),
			markEmpty(trimPart(inf.uuid())),
			markEmpty(trimPart(inf.auth())),
			markEmpty(inf.host()),
			markEmpty(inf.method()),
			markEmpty(inf.path()),
			bytefmt.ByteSize(uint64(inf.clen())),
			bytefmt.ByteSize(uint64(inf.size())),
			markEmpty(inf.time()),
			markEmpty(inf.user()),
			inf.fail(),
		)
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

func trimPart(s string) string {
	if len(s) > magicLen {
		return s[:magicLen]
	}
	return s
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

func (i informer) host() string {
	return hostFromCtx(i.c)
}

func (i informer) user() string {
	return fmt.Sprintf("%q", userFromCtx(i.c))
}

func (i informer) uuid() string {
	return uuidFromCtx(i.c)[:magicLen]
}

func (i informer) auth() string {
	return authFromCtx(i.c)
}

func (i informer) code() int64 {
	return codeFromCtx(i.c)
}

func (i informer) clen() int64 {
	return clenFromCtx(i.c)
}

func (i informer) size() int64 {
	return sizeFromCtx(i.c)
}

func (i informer) fail() string {
	err := failFromCtx(i.c)
	if err != nil {
		return fmt.Sprintf("err: %v", err)
	}
	return ""
}

func (i informer) time() string {
	if t := timeFromCtx(i.c); !t.IsZero() {
		return time.Since(t).String()
	}
	return ""
}
