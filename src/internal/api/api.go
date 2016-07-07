package api

import (
	"fmt"
	"net/http"
	"strings"

	"internal/ctxutil"
	"internal/pipe"
	"internal/version"
)

type headReader interface {
	ReadHeader(http.Header)
}

type headWriter interface {
	WriteHeader(http.Header)
}

type newer interface {
	New() interface{}
}

type worker interface {
	Work([]byte) (interface{}, error)
}

type workFunc func([]byte) (interface{}, error)

func (f workFunc) Work(b []byte) (interface{}, error) {
	return f(b)
}

var (
	httpHandlers = map[string]http.Handler{
		"GET>/":     pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(root), pipe.Resp, pipe.Tail),
		"GET>/ping": pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(work(workFunc(Ping))), pipe.Resp, pipe.Tail),

		"POST>/system/get-auth": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(GetAuth))), pipe.Resp, pipe.Tail),
		"POST>/system/set-auth": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(SetAuth))), pipe.Resp, pipe.Tail),
		"POST>/system/del-auth": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(DelAuth))), pipe.Resp, pipe.Tail),

		"POST>/system/get-addr": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(GetAddr))), pipe.Resp, pipe.Tail),
		"POST>/system/set-addr": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(SetAddr))), pipe.Resp, pipe.Tail),
		"POST>/system/del-addr": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(DelAddr))), pipe.Resp, pipe.Tail),

		"POST>/system/get-drug": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(GetDrug))), pipe.Resp, pipe.Tail),
		"POST>/system/set-drug": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(SetDrug))), pipe.Resp, pipe.Tail),
		"POST>/system/del-drug": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(DelDrug))), pipe.Resp, pipe.Tail),

		"POST>/system/get-stat": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(GetStat))), pipe.Resp, pipe.Tail),
		"POST>/system/set-stat": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(SetStat))), pipe.Resp, pipe.Tail),
		"POST>/system/del-stat": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(DelStat))), pipe.Resp, pipe.Tail),

		"POST>/system/get-meta": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(GetMeta))), pipe.Resp, pipe.Tail),
		"POST>/system/get-zlog": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(workFunc(GetZlog))), pipe.Resp, pipe.Tail),

		// Converter from old school style /data/add DEPRECATED
		"POST>/data/add": pipe.Use(pipe.Conv, pipe.Head, pipe.Auth(0), pipe.Meta, pipe.Work(work(Putd)), pipe.Resp, pipe.Tail),

		"POST>/stream/put-data": pipe.Use(pipe.Head, pipe.Auth(0), pipe.Meta, pipe.Work(work(Putd)), pipe.Resp, pipe.Tail),
		"POST>/stream/pop-data": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(Popd)), pipe.Resp, pipe.Tail),
		"POST>/stream/get-data": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(Getd)), pipe.Resp, pipe.Tail),
		"POST>/stream/del-data": pipe.Use(pipe.Head, pipe.Auth(1), pipe.Gzip, pipe.Work(work(Deld)), pipe.Resp, pipe.Tail),

		// => Debug mode only, when flag.Debug == true
		"GET>/debug/info":               pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(work(workFunc(Info))), pipe.Resp, pipe.Tail), // ?
		"GET>/debug/vars":               pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(stdh), pipe.Resp, pipe.Tail),                 // expvar
		"GET>/debug/pprof/":             pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(stdh), pipe.Resp, pipe.Tail),                 // net/http/pprof
		"GET>/debug/pprof/cmdline":      pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(stdh), pipe.Resp, pipe.Tail),                 // net/http/pprof
		"GET>/debug/pprof/profile":      pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(stdh), pipe.Resp, pipe.Tail),                 // net/http/pprof
		"GET>/debug/pprof/symbol":       pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(stdh), pipe.Resp, pipe.Tail),                 // net/http/pprof
		"GET>/debug/pprof/trace":        pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(stdh), pipe.Resp, pipe.Tail),                 // net/http/pprof
		"GET>/debug/pprof/goroutine":    pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(stdh), pipe.Resp, pipe.Tail),                 // runtime/pprof
		"GET>/debug/pprof/threadcreate": pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(stdh), pipe.Resp, pipe.Tail),                 // runtime/pprof
		"GET>/debug/pprof/heap":         pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(stdh), pipe.Resp, pipe.Tail),                 // runtime/pprof
		"GET>/debug/pprof/block":        pipe.Use(pipe.Head, pipe.Gzip, pipe.Work(stdh), pipe.Resp, pipe.Tail),                 // runtime/pprof

		// => Workarounds for 404/405
		"GET>/error/404": pipe.Use(pipe.Head, pipe.Work(resp404), pipe.Resp, pipe.Tail),
		"GET>/error/405": pipe.Use(pipe.Head, pipe.Work(resp405), pipe.Resp, pipe.Tail),
	}
)

func root(w http.ResponseWriter, r *http.Request) {
	*r = *r.WithContext(
		ctxutil.WithData(
			r.Context(),
			fmt.Sprintf("%s %s", version.AppName(), version.WithBuildInfo())),
	)
}

func respErr(r *http.Request, code int) {
	*r = *r.WithContext(
		ctxutil.WithFail(
			r.Context(),
			fmt.Errorf("api: %s", strings.ToLower(http.StatusText(code))), code),
	)
}

func resp404(w http.ResponseWriter, r *http.Request) {
	respErr(r, http.StatusNotFound)
}

func resp405(w http.ResponseWriter, r *http.Request) {
	respErr(r, http.StatusMethodNotAllowed)
}
