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

type handlerPipe func(h http.HandlerFunc) http.HandlerFunc
type bundle struct {
	h http.Handler
	w core.Worker
}

var (
	mapCoreWorkers  map[string]core.Worker
	mapHTTPHandlers = map[string]bundle{
		"GET:/":     {use(pipeHead, pipeGzip, pipe(root), pipeResp, pipeTail), nil},
		"GET:/ping": {use(pipeHead, pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.Ping)},

		"POST:/system/get-auth": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.GetAuth)},
		"POST:/system/set-auth": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.SetAuth)},
		"POST:/system/del-auth": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.DelAuth)},

		"POST:/system/get-addr": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.GetAddr)},
		"POST:/system/set-addr": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.SetAddr)},
		"POST:/system/del-addr": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.DelAddr)},

		"POST:/system/get-drug": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.GetDrug)},
		"POST:/system/set-drug": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.SetDrug)},
		"POST:/system/del-drug": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.DelDrug)},

		"POST:/system/get-stat": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.GetStat)},
		"POST:/system/set-stat": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.SetStat)},
		"POST:/system/del-stat": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.DelStat)},

		// Converter from old school style data/add DEPRECATED
		"POST:/data/add": {use(pipeConv, pipeHead, pipeAuth(0), pipeMeta, pipe(work), pipeResp, pipeTail), core.Putd},

		"POST:/stream/put-data": {use(pipeHead, pipeAuth(0), pipeMeta, pipe(work), pipeResp, pipeTail), core.Putd},
		"POST:/stream/pop-data": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.Popd},
		"POST:/stream/get-data": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.Getd},
		"POST:/stream/del-data": {use(pipeHead, pipeAuth(1), pipeGzip, pipe(work), pipeResp, pipeTail), core.Deld},

		// => Debug mode only, when flag.Debug == true
		"GET:/debug/info":               {use(pipeHead, pipeGzip, pipe(work), pipeResp, pipeTail), core.WorkFunc(core.Info)}, // ?
		"GET:/debug/vars":               {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                      // expvar
		"GET:/debug/pprof/":             {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                      // net/http/pprof
		"GET:/debug/pprof/cmdline":      {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                      // net/http/pprof
		"GET:/debug/pprof/profile":      {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                      // net/http/pprof
		"GET:/debug/pprof/symbol":       {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                      // net/http/pprof
		"GET:/debug/pprof/trace":        {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                      // net/http/pprof
		"GET:/debug/pprof/goroutine":    {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                      // runtime/pprof
		"GET:/debug/pprof/threadcreate": {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                      // runtime/pprof
		"GET:/debug/pprof/heap":         {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                      // runtime/pprof
		"GET:/debug/pprof/block":        {use(pipeHead, pipeGzip, pipe(stdh), pipeResp, pipeTail), nil},                      // runtime/pprof

		// => Workarounds for 404/405
		"GET:/error/404": {use(pipeHead, pipe(resp404), pipeResp, pipeTail), nil},
		"GET:/error/405": {use(pipeHead, pipe(resp405), pipeResp, pipeTail), nil},
	}
)

// Reg is caled from main package for manually initialization
func Reg(reg func(m, p string, h http.Handler)) error {
	mapCoreWorkers = make(map[string]core.Worker, len(mapHTTPHandlers))
	for k, v := range mapHTTPHandlers {
		s := strings.Split(k, ":")
		if reg != nil {
			reg(s[0], s[1], v.h)
		}
		if v.w != nil {
			mapCoreWorkers[s[1]] = v.w
		}
	}

	return core.Init()
}

func root(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = ctxWithData(ctx, fmt.Sprintf("%s %s", version.AppName(), version.WithBuildInfo()))
	*r = *r.WithContext(ctx)
}

func work(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	ctx := r.Context()
	wrk, ok := mapCoreWorkers[r.URL.Path]
	if err := fmt.Errorf("api: core method not found"); !ok {
		*r = *r.WithContext(ctxWithFail(ctx, err))
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

	if m, ok := wrk.(core.Master); ok {
		wrk = m.NewWorker()
	}

	// 1
	if hr, ok := wrk.(core.HTTPHeadReader); ok {
		hr.ReadHeader(r.Header)
	}

	// 2
	out, err := wrk.Work(buf.Bytes())
	if err != nil {
		*r = *r.WithContext(ctxWithFail(ctx, err))
		return
	}

	// 3
	if hw, ok := wrk.(core.HTTPHeadWriter); ok {
		hw.WriteHeader(w.Header())
	}

	ctx = ctxWithData(ctx, out)
	*r = *r.WithContext(ctx)
}

func respErr(r *http.Request, code int) {
	ctx := r.Context()
	ctx = ctxWithFail(ctx, fmt.Errorf("api: %s", strings.ToLower(http.StatusText(code))), code)
	*r = *r.WithContext(ctx)
}

func resp404(w http.ResponseWriter, r *http.Request) {
	respErr(r, http.StatusNotFound)
}

func resp405(w http.ResponseWriter, r *http.Request) {
	respErr(r, http.StatusMethodNotAllowed)
}
