package api

import (
	"bytes"
	"fmt"
	"io"
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

type handlerPair struct {
	h http.Handler
	w worker
}

var (
	mapCoreWorkers  map[string]worker
	mapHTTPHandlers = map[string]handlerPair{
		"GET:/":     {use(pipeHead, pipeGzip, pipe(root), pipeResp, pipeTail), nil},
		"GET:/ping": {use(pipeHead, pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.Ping)},

		"POST:/system/get-auth": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.GetAuth)},
		"POST:/system/set-auth": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.SetAuth)},
		"POST:/system/del-auth": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.DelAuth)},

		"POST:/system/get-addr": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.GetAddr)},
		"POST:/system/set-addr": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.SetAddr)},
		"POST:/system/del-addr": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.DelAddr)},

		"POST:/system/get-drug": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.GetDrug)},
		"POST:/system/set-drug": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.SetDrug)},
		"POST:/system/del-drug": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.DelDrug)},

		"POST:/system/get-stat": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.GetStat)},
		"POST:/system/set-stat": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.SetStat)},
		"POST:/system/del-stat": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.DelStat)},

		// Converter from old school style /data/add DEPRECATED
		"POST:/data/add": {use(pipeConv, pipeHead, pipeAuth(0), pipeMeta, pipe(work), pipeResp, pipeTail), core.Putd},

		"POST:/stream/put-data": {use(pipeHead, pipeAuth(0), pipeMeta, pipe(work), pipeResp, pipeTail), core.Putd},
		"POST:/stream/pop-data": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.Popd},
		"POST:/stream/get-data": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.Getd},
		"POST:/stream/del-data": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.Deld},

		// => Debug mode only, when flag.Debug == true
		"GET:/debug/info":               {use(pipeHead, pipeGzip, pipe(work), pipeResp, pipeTail), workFunc(core.Info)}, // ?
		"GET:/debug/vars":               {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                 // expvar
		"GET:/debug/pprof/":             {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                 // net/http/pprof
		"GET:/debug/pprof/cmdline":      {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                 // net/http/pprof
		"GET:/debug/pprof/profile":      {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                 // net/http/pprof
		"GET:/debug/pprof/symbol":       {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                 // net/http/pprof
		"GET:/debug/pprof/trace":        {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                 // net/http/pprof
		"GET:/debug/pprof/goroutine":    {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                 // runtime/pprof
		"GET:/debug/pprof/threadcreate": {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                 // runtime/pprof
		"GET:/debug/pprof/heap":         {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                 // runtime/pprof
		"GET:/debug/pprof/block":        {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                 // runtime/pprof

		// => Workarounds for 404/405
		"GET:/error/404": {use(pipeHead, pipe(resp404), pipeResp, pipeTail), nil},
		"GET:/error/405": {use(pipeHead, pipe(resp405), pipeResp, pipeTail), nil},
	}
)

// Reg is called manually initialization
func Reg(reg func(m, p string, h http.Handler)) error {
	mapCoreWorkers = make(map[string]worker, len(mapHTTPHandlers))
	for k, v := range mapHTTPHandlers {
		s := strings.Split(k, ":")
		if reg != nil {
			reg(s[0], s[1], v.h)
		}
		if v.w != nil {
			mapCoreWorkers[s[1]] = v.w
		}
	}

	return nil //core.Init()
}

func root(w http.ResponseWriter, r *http.Request) {
	*r = *r.WithContext(
		ctxWithData(
			r.Context(),
			fmt.Sprintf("%s %s", version.AppName(), version.WithBuildInfo())),
	)
}

func work(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	ctx := r.Context()
	wrk, ok := mapCoreWorkers[r.URL.Path]
	if !ok {
		*r = *r.WithContext(ctxWithFail(ctx, fmt.Errorf("api: core method not found")))
		return
	}

	var buf = new(bytes.Buffer)
	if r.Method == "POST" {
		n, err := io.Copy(buf, r.Body)
		if err != nil {
			*r = *r.WithContext(ctxWithFail(ctx, err))
			return
		}
		ctx = ctxWithClen(ctx, n)
	}

	if nwr, ok := wrk.(newer); ok {
		if nwr, ok := nwr.(worker); ok {
			wrk = nwr
		}
	}

	// 1
	if hr, ok := wrk.(headReader); ok {
		hr.ReadHeader(r.Header)
	}

	// 2
	out, err := wrk.Work(buf.Bytes())
	if err != nil {
		*r = *r.WithContext(ctxWithFail(ctx, err))
		return
	}

	// 3
	if hw, ok := wrk.(headWriter); ok {
		hw.WriteHeader(w.Header())
	}

	ctx = ctxWithData(ctx, out)
	*r = *r.WithContext(ctx)
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
