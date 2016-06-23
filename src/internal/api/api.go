package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"

	"internal/core"
	"internal/pref"
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
	s := fmt.Sprintf("%s %s", version.AppName(), version.WithBuildInfo())

	r ctxWith200(w, r, s)
}

func work(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	wrk, ok := mapCoreWorkers[r.URL.Path]
	if !ok {
		return resp500(ctx, fmt.Errorf("api: core method not found"))
	}

	var buf = new(bytes.Buffer)
	if r.Method == "POST" {
		n, err := io.Copy(buf, r.Body)
		if err != nil {
			return with500(ctxWithClen(ctx, n), err)
		}
		ctx = ctxWithClen(ctx, n)
	}

	if m, ok := wrk.(core.Master); ok {
		wrk = m.NewWorker()
	}

	if hr, ok := wrk.(core.HTTPHeadReader); ok {
		hr.ReadHeader(r.Header)
	}

	res, err := wrk.Work(buf.Bytes())
	if err != nil {
		return with500(ctx, err)
	}

	if hw, ok := wrk.(core.HTTPHeadWriter); ok {
		hw.WriteHeader(w.Header())
	}

	return resp200(ctx, w, res)
}

func stdh(w http.ResponseWriter, r *http.Request) {
	if !pref.Debug {
		return resp500(ctx, fmt.Errorf("api: flag debug not found"))
	}

	if h, p := http.DefaultServeMux.Handler(r); p != "" {
		h.ServeHTTP(w, r)
		return ctxWithCode(ctxWithSize(ctx, 0), http.StatusOK) // TODO: wrap w to get real size
	}

	return ctxWithSize(ctxWithFail(ctx, fmt.Errorf("api: unreachable")), 0)
}

func resp404(w http.ResponseWriter, r *http.Request) {
	return respErr(ctx, http.StatusNotFound)
}

func resp405(w http.ResponseWriter, r *http.Request) {
	return respErr(ctx, http.StatusMethodNotAllowed)
}
