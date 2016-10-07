package api

import (
	"fmt"
	"net/http"
	"strings"

	"internal/core"
	"internal/core/ctxt"
	"internal/core/pipe"
	"internal/core/pref"
	"internal/version"

	"github.com/julienschmidt/httprouter"
)

var (
	httpHandlers = map[string]http.Handler{
		"GET>/":     pipe.Use(pipe.Head, pipe.Gzip, pipe.Wrap(root), pipe.Resp, pipe.Tail),
		"GET>/ping": pipe.Use(pipe.Head, pipe.Gzip, pipe.Wrap(ping), pipe.Resp, pipe.Tail),

		"POST>/system/get-auth": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(getAuth), pipe.Resp, pipe.Tail),
		"POST>/system/set-auth": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(setAuth), pipe.Resp, pipe.Tail),
		"POST>/system/del-auth": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(delAuth), pipe.Resp, pipe.Tail),

		"POST>/system/get-addr": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(getAddr), pipe.Resp, pipe.Tail),
		"POST>/system/set-addr": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(setAddr), pipe.Resp, pipe.Tail),
		"POST>/system/del-addr": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(delAddr), pipe.Resp, pipe.Tail),

		"POST>/system/get-drug": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(getDrug), pipe.Resp, pipe.Tail),
		"POST>/system/set-drug": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(setDrug), pipe.Resp, pipe.Tail),
		"POST>/system/del-drug": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(delDrug), pipe.Resp, pipe.Tail),

		"POST>/system/get-stat": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(getStat), pipe.Resp, pipe.Tail),
		"POST>/system/set-stat": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(setStat), pipe.Resp, pipe.Tail),
		"POST>/system/del-stat": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(delStat), pipe.Resp, pipe.Tail),

		"POST>/system/get-meta": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(getMeta), pipe.Resp, pipe.Tail),
		"POST>/system/get-zlog": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(getZlog), pipe.Resp, pipe.Tail),

		// Converter from old school style /data/add DEPRECATED
		"POST>/data/add": pipe.Use(pipe.Conv, pipe.Head, pipe.Auth(pass), pipe.Meta, pipe.Wrap(putd), pipe.Resp, pipe.Tail),

		"POST>/stream/put-data": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Meta, pipe.Wrap(putd), pipe.Resp, pipe.Tail),
		"POST>/stream/pop-data": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(popd), pipe.Resp, pipe.Tail),
		"POST>/stream/get-data": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(getd), pipe.Resp, pipe.Tail),
		"POST>/stream/del-data": pipe.Use(pipe.Head, pipe.Auth(pass), pipe.Gzip, pipe.Wrap(deld), pipe.Resp, pipe.Tail),

		// => Debug mode only, when pref.Debug == true
		"GET>/debug/info":               pipe.Use(pipe.Head, pipe.Gzip, pipe.Wrap(info), pipe.Resp, pipe.Tail), // ?
		"GET>/debug/vars":               pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // expvar
		"GET>/debug/pprof/":             pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // net/http/pprof
		"GET>/debug/pprof/cmdline":      pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // net/http/pprof
		"GET>/debug/pprof/profile":      pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // net/http/pprof
		"GET>/debug/pprof/symbol":       pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // net/http/pprof
		"GET>/debug/pprof/trace":        pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // net/http/pprof
		"GET>/debug/pprof/goroutine":    pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // runtime/pprof
		"GET>/debug/pprof/threadcreate": pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // runtime/pprof
		"GET>/debug/pprof/heap":         pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // runtime/pprof
		"GET>/debug/pprof/block":        pipe.Use(pipe.Head, pipe.Gzip, pipe.StdH, pipe.Resp, pipe.Tail),       // runtime/pprof

		// => Workarounds for 404/405
		"GET>/error/404": pipe.Use(pipe.Head, pipe.Wrap(e404), pipe.Resp, pipe.Tail),
		"GET>/error/405": pipe.Use(pipe.Head, pipe.Wrap(e405), pipe.Resp, pipe.Tail),
	}
)

func pass(key string) bool {
	if strings.EqualFold(pref.MasterKey, key) {
		return true
	}
	return core.Pass(key)
}

func root() (interface{}, error) {
	return fmt.Sprintf("%s %s", version.AppName(), version.WithBuildInfo()), nil
}

func e404(w http.ResponseWriter, r *http.Request) {
	respErr(r, http.StatusNotFound)
}

func e405(w http.ResponseWriter, r *http.Request) {
	respErr(r, http.StatusMethodNotAllowed)
}

func respErr(r *http.Request, code int) {
	ctx := r.Context()
	err := fmt.Errorf("api: %s", strings.ToLower(http.StatusText(code)))
	ctx = ctxt.WithFail(ctx, err, code)
	*r = *r.WithContext(ctx)
}

// MakeRouter returns http.Handler
func MakeRouter() http.Handler {
	r := httprouter.New()

	for k, v := range httpHandlers {
		s := strings.Split(k, ">") // [m,p]

		switch s[1] {
		case "/error/404":
			r.NotFound = v
		case "/error/405":
			r.MethodNotAllowed = v
		default:
			func(m, p string, h http.Handler) {
				r.Handle(m, p,
					func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
						ctx := r.Context()
						for i := range p {
							ctx = ctxt.WithURLp(ctx, p[i].Key, p[i].Value)
						}
						r = r.WithContext(ctx)
						h.ServeHTTP(w, r)
					})
			}(s[0], s[1], v)
		}
	}

	return r
}
