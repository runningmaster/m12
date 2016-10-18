package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"internal/core"
	"internal/core/pipe"
	"internal/version"

	"github.com/julienschmidt/httprouter"
)

const (
	code404 = http.StatusNotFound
	code405 = http.StatusMethodNotAllowed
)

var (
	fake404URL = fmt.Sprintf("/workaround-%d", code404)
	fake405URL = fmt.Sprintf("/workaround-%d", code404)

	httpHandlers = map[string]http.Handler{
		"GET /":     pipe.Use(pipe.Head, pipe.Gzip, pipe.Wrap(root), pipe.Resp, pipe.Tail),
		"GET /ping": pipe.Use(pipe.Head, pipe.Gzip, pipe.Wrap(ping), pipe.Resp, pipe.Tail),

		"POST /system/get-auth": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(core.GetAuth), pipe.Resp, pipe.Tail),
		"POST /system/set-auth": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(core.SetAuth), pipe.Resp, pipe.Tail),
		"POST /system/del-auth": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(core.DelAuth), pipe.Resp, pipe.Tail),

		"POST /system/get-addr": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(core.GetAddr), pipe.Resp, pipe.Tail),
		"POST /system/set-addr": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(core.SetAddr), pipe.Resp, pipe.Tail),
		"POST /system/del-addr": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(core.DelAddr), pipe.Resp, pipe.Tail),

		"POST /system/get-drug": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(core.GetDrug), pipe.Resp, pipe.Tail),
		"POST /system/set-drug": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(core.SetDrug), pipe.Resp, pipe.Tail),
		"POST /system/del-drug": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(core.DelDrug), pipe.Resp, pipe.Tail),

		"POST /system/get-stat": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(core.GetStat), pipe.Resp, pipe.Tail),
		"POST /system/set-stat": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(core.SetStat), pipe.Resp, pipe.Tail),
		"POST /system/del-stat": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(core.DelStat), pipe.Resp, pipe.Tail),

		"POST /system/get-meta": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(getMeta), pipe.Resp, pipe.Tail),
		"POST /system/get-zlog": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(getZlog), pipe.Resp, pipe.Tail),

		// DEPRECATED Converter from old school style /data/add
		"POST /data/add": pipe.Use(pipe.Conv, pipe.Head, pipe.Auth(pass), pipe.Meta, pipe.Wrap(putd), pipe.Resp, pipe.Tail),

		"POST /stream/put-data": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Meta, pipe.Wrap(putd), pipe.Resp, pipe.Tail),
		"POST /stream/pop-data": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(popd), pipe.Resp, pipe.Tail),
		"POST /stream/get-data": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(getd), pipe.Resp, pipe.Tail),
		"POST /stream/del-data": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(deld), pipe.Resp, pipe.Tail),

		// => Debug mode only, when pref.Debug == true
		"GET /debug/vars":               pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // expvar
		"GET /debug/pprof/":             pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // net/http/pprof
		"GET /debug/pprof/cmdline":      pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // net/http/pprof
		"GET /debug/pprof/profile":      pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // net/http/pprof
		"GET /debug/pprof/symbol":       pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // net/http/pprof
		"GET /debug/pprof/trace":        pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // net/http/pprof
		"GET /debug/pprof/goroutine":    pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // runtime/pprof
		"GET /debug/pprof/threadcreate": pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // runtime/pprof
		"GET /debug/pprof/heap":         pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // runtime/pprof
		"GET /debug/pprof/block":        pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // runtime/pprof
		"GET /debug/redis/info":         pipe.Use(pipe.Head, pipe.Gzip, pipe.Wrap(info), pipe.Resp, pipe.Tail), // ?

		// => Workarounds for 404/405
		"GET " + fake404URL: pipe.Use(pipe.Head, pipe.ErrH(code404), pipe.Resp, pipe.Tail),
		"GET " + fake405URL: pipe.Use(pipe.Head, pipe.ErrH(code405), pipe.Resp, pipe.Tail),
	}
)

// MakeRouter returns http.Handler
func MakeRouter() http.Handler {
	r := httprouter.New()

	for k, v := range httpHandlers {
		s := strings.Split(k, " ") // [m,p]

		switch s[1] {
		case fake404URL:
			r.NotFound = v
		case fake405URL:
			r.MethodNotAllowed = v
		default:
			func(m, p string, h http.Handler) {
				r.Handle(m, p,
					func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
						ctx := r.Context()
						for i := range p {
							ctx = context.WithValue(ctx, p[i].Key, p[i].Value)
						}
						r = r.WithContext(ctx)
						h.ServeHTTP(w, r)
					})
			}(s[0], s[1], v)
		}
	}

	return r
}

func pass(key string) bool {
	return core.Pass(key)
}

func root() (interface{}, error) {
	return fmt.Sprintf("%s %s", version.AppName(), version.WithBuildInfo()), nil
}

func ping() (interface{}, error) {
	return core.Ping()
}

func info() (interface{}, error) {
	return core.Info()
}
