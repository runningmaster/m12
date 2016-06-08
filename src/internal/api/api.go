package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"

	"internal/core"
	"internal/pref"
	"internal/server"
	"internal/version"

	"golang.org/x/net/context"
)

type handlerFunc func(context.Context, http.ResponseWriter, *http.Request)
type handlerFuncCtx func(context.Context, http.ResponseWriter, *http.Request) context.Context
type handlerPipe func(h handlerFunc) handlerFunc
type bundle struct {
	h http.Handler
	w core.Worker
}

func (f handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(context.Background(), w, r)
}

var (
	mapCoreWorkers  map[string]core.Worker
	mapHTTPHandlers = map[string]bundle{
		"GET:/":     {use(pipeHead, pipeGzip, pipe(root), pipeFail, pipeTail), nil},
		"GET:/ping": {use(pipeHead, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.Ping)},

		"POST:/system/get-auth": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.GetAuth)},
		"POST:/system/set-auth": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.SetAuth)},
		"POST:/system/del-auth": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.DelAuth)},

		"POST:/system/get-link-addr": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.GetLinkAddr)},
		"POST:/system/set-link-addr": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.SetLinkAddr)},
		"POST:/system/del-link-addr": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.DelLinkAddr)},

		"POST:/system/get-link-drug": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.GetLinkDrug)},
		"POST:/system/set-link-drug": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.SetLinkDrug)},
		"POST:/system/del-link-drug": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.DelLinkDrug)},

		"POST:/system/get-link-stat": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.GetLinkStat)},
		"POST:/system/set-link-stat": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.SetLinkStat)},
		"POST:/system/del-link-stat": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.DelLinkStat)},

		"POST:/stream/put-data": {use(pipeHead, pipeAuth, pipeMeta, pipe(exec), pipeFail, pipeTail), core.Putd},
		"POST:/stream/pop-data": {use(pipeHead, pipeAuth, pipeGzip, pipe(exec), pipeFail, pipeTail), core.Popd},

		// Converter from old school style data/add DEPRECATED
		"POST:/data/add": {use(pipeConv, pipeHead, pipeAuth, pipeMeta, pipe(exec), pipeFail, pipeTail), core.Conv},

		// => Debug mode only, when flag.Debug == true
		"GET:/debug/info":               {use(pipeHead, pipeGzip, pipe(exec), pipeFail, pipeTail), core.WorkFunc(core.Info)}, // ?
		"GET:/debug/vars":               {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},                      // expvar
		"GET:/debug/pprof/":             {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},                      // net/http/pprof
		"GET:/debug/pprof/cmdline":      {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},                      // net/http/pprof
		"GET:/debug/pprof/profile":      {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},                      // net/http/pprof
		"GET:/debug/pprof/symbol":       {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},                      // net/http/pprof
		"GET:/debug/pprof/trace":        {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},                      // net/http/pprof
		"GET:/debug/pprof/goroutine":    {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},                      // runtime/pprof
		"GET:/debug/pprof/threadcreate": {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},                      // runtime/pprof
		"GET:/debug/pprof/heap":         {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},                      // runtime/pprof
		"GET:/debug/pprof/block":        {use(pipeHead, pipeGzip, pipe(stdh), pipeFail, pipeTail), nil},                      // runtime/pprof

		// => Workarounds for 404/405
		"GET:/error/404": {use(pipeHead, pipe(with404), pipeFail, pipeTail), nil},
		"GET:/error/405": {use(pipeHead, pipe(with405), pipeFail, pipeTail), nil},
	}
)

// Reg is caled from main package for manually initialization
func Reg() error {
	mapCoreWorkers = make(map[string]core.Worker, len(mapHTTPHandlers))
	for k, v := range mapHTTPHandlers {
		s := strings.Split(k, ":")
		server.RegHandler(s[0], s[1], v.h)
		if v.w != nil {
			mapCoreWorkers[s[1]] = v.w
		}
	}

	return core.Init()
}

func root(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	res := fmt.Sprintf("%s %s", version.AppName(), version.WithBuildInfo())
	return with200(ctx, w, res)
}

func exec(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	defer func() { _ = r.Body.Close() }()

	wrk, ok := mapCoreWorkers[r.URL.Path]
	if !ok {
		return with500(ctx, fmt.Errorf("api: core method not found"))
	}

	var buf = new(bytes.Buffer)
	if r.Method == "POST" {
		n, err := io.Copy(buf, r.Body)
		if err != nil {
			return with500(ctxWithClen(ctx, n), err)
		}
		ctx = ctxWithClen(ctx, n)
	}

	var hr core.HTTPHeadReader
	if hr, ok = wrk.(core.HTTPHeadReader); ok {
		hr.ReadHeader(r.Header)
	}

	res, err := wrk.Work(buf.Bytes())
	if err != nil {
		return with500(ctx, err)
	}

	var hw core.HTTPHeadWriter
	if hw, ok = wrk.(core.HTTPHeadWriter); ok {
		hw.WriteHeader(w.Header())
	}

	return with200(ctx, w, res)
}

func stdh(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	if !pref.Debug {
		return with500(ctx, fmt.Errorf("api: flag debug not found"))
	}

	if h, p := http.DefaultServeMux.Handler(r); p != "" {
		h.ServeHTTP(w, r)
		return ctxWithCode(ctxWithSize(ctx, 0), http.StatusOK) // TODO: wrap w to get real size
	}

	return ctxWithSize(ctxWithFail(ctx, fmt.Errorf("api: unreachable")), 0)
}

func writeJSON(ctx context.Context, w http.ResponseWriter, code int, data interface{}) (int64, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Powered-By", runtime.Version())
	w.Header().Set("X-Request-ID", uuidFromCtx(ctx))
	w.WriteHeader(code)

	if true { // FIXME (flag?)
		var tmp bytes.Buffer
		err = json.Indent(&tmp, b, "", "\t")
		if err != nil {
			return 0, err
		}
		b = tmp.Bytes()
	}

	n, err := w.Write(b)
	if err != nil {
		return 0, err
	}

	_, err = w.Write([]byte("\n"))
	if err != nil {
		return 0, err
	}
	n++

	return int64(n), nil
}

func with200(ctx context.Context, w http.ResponseWriter, res interface{}) context.Context {
	n, err := writeJSON(ctx, w, http.StatusOK, res)
	if err != nil {
		return ctxWithSize(ctxWithFail(ctx, err), n)
	}
	return ctxWithCode(ctxWithSize(ctx, n), http.StatusOK)
}

func with500(ctx context.Context, err error) context.Context {
	return ctxWithCode(ctxWithFail(ctx, err), http.StatusInternalServerError)
}

func withErr(ctx context.Context, code int) context.Context {
	err := fmt.Errorf("api: %s", strings.ToLower(http.StatusText(code)))
	return ctxWithCode(ctxWithFail(ctx, err), int64(code))
}

func with404(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return withErr(ctx, http.StatusNotFound)
}

func with405(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return withErr(ctx, http.StatusMethodNotAllowed)
}
