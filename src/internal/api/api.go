package api

import (
	"fmt"

	"net/http"
	"strings"

	"internal/core"
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

type handlerPipe func(h http.Handler) http.Handler

var (
	httpHandlers = map[string]http.Handler{
		"GET>/":     use(pipeHead, pipeGzip, pipe(root), pipeResp, pipeTail),
		"GET>/ping": use(pipeHead, pipeGzip, pipe(work(workFunc(core.Ping))), pipeResp, pipeTail),

		"POST>/system/get-auth": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.GetAuth))), pipeResp, pipeTail),
		"POST>/system/set-auth": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.SetAuth))), pipeResp, pipeTail),
		"POST>/system/del-auth": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.DelAuth))), pipeResp, pipeTail),

		"POST>/system/get-addr": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.GetAddr))), pipeResp, pipeTail),
		"POST>/system/set-addr": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.SetAddr))), pipeResp, pipeTail),
		"POST>/system/del-addr": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.DelAddr))), pipeResp, pipeTail),

		"POST>/system/get-drug": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.GetDrug))), pipeResp, pipeTail),
		"POST>/system/set-drug": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.SetDrug))), pipeResp, pipeTail),
		"POST>/system/del-drug": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.DelDrug))), pipeResp, pipeTail),

		"POST>/system/get-stat": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.GetStat))), pipeResp, pipeTail),
		"POST>/system/set-stat": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.SetStat))), pipeResp, pipeTail),
		"POST>/system/del-stat": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.DelStat))), pipeResp, pipeTail),

		"POST>/system/get-meta": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.GetMeta))), pipeResp, pipeTail),
		"POST>/system/get-zlog": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(workFunc(core.GetZlog))), pipeResp, pipeTail),

		// Converter from old school style /data/add DEPRECATED
		"POST>/data/add": use(pipeConv, pipeHead, pipeAuth(0), pipeMeta, pipe(work(core.Putd)), pipeResp, pipeTail),

		"POST>/stream/put-data": use(pipeHead, pipeAuth(0), pipeMeta, pipe(work(core.Putd)), pipeResp, pipeTail),
		"POST>/stream/pop-data": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(core.Popd)), pipeResp, pipeTail),
		"POST>/stream/get-data": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(core.Getd)), pipeResp, pipeTail),
		"POST>/stream/del-data": use(pipeHead, pipeAuth(1), pipeGzip, pipe(work(core.Deld)), pipeResp, pipeTail),

		// => Debug mode only, when flag.Debug == true
		"GET>/debug/info":               use(pipeHead, pipeGzip, pipe(work(workFunc(core.Info))), pipeResp, pipeTail), // ?
		"GET>/debug/vars":               use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail),                      // expvar
		"GET>/debug/pprof/":             use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail),                      // net/http/pprof
		"GET>/debug/pprof/cmdline":      use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail),                      // net/http/pprof
		"GET>/debug/pprof/profile":      use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail),                      // net/http/pprof
		"GET>/debug/pprof/symbol":       use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail),                      // net/http/pprof
		"GET>/debug/pprof/trace":        use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail),                      // net/http/pprof
		"GET>/debug/pprof/goroutine":    use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail),                      // runtime/pprof
		"GET>/debug/pprof/threadcreate": use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail),                      // runtime/pprof
		"GET>/debug/pprof/heap":         use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail),                      // runtime/pprof
		"GET>/debug/pprof/block":        use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail),                      // runtime/pprof

		// => Workarounds for 404/405
		"GET>/error/404": use(pipeHead, pipe(resp404), pipeResp, pipeTail),
		"GET>/error/405": use(pipeHead, pipe(resp405), pipeResp, pipeTail),
	}
)

func root(w http.ResponseWriter, r *http.Request) {
	*r = *r.WithContext(
		ctxWithData(
			r.Context(),
			fmt.Sprintf("%s %s", version.AppName(), version.WithBuildInfo())),
	)
}

func respErr(r *http.Request, code int) {
	*r = *r.WithContext(
		ctxWithFail(
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
