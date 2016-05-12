package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"internal/conf"
	"internal/core"
	"internal/version"

	"golang.org/x/net/context"
)

type handlerFunc func(context.Context, http.ResponseWriter, *http.Request)
type handlerFuncCtx func(context.Context, http.ResponseWriter, *http.Request) context.Context
type handlerPipe func(h handlerFunc) handlerFunc
type bundle struct {
	h http.Handler
	f core.Handler
}

func (f handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(context.Background(), w, r)
}

func root(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	res := fmt.Sprintf("%s %s", version.AppName(), version.String())
	return with200(ctx, w, res)
}

func exec(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	f, ok := mapCoreHandlers[r.URL.Path]
	if !ok {
		return with500(ctx, fmt.Errorf("api: core method not found"))
	}

	res, err := f(ctx, w, r)
	if err != nil {
		return with500(ctx, err)
	}

	return with200(ctx, w, res)
}

func stdh(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	if !conf.Debug {
		return with500(ctx, fmt.Errorf("api: flag debug not found"))
	}

	if h, p := http.DefaultServeMux.Handler(r); p != "" {
		h.ServeHTTP(w, r)
		return withCode(withSize(ctx, 0), http.StatusOK) // TODO: wrap w to get real size
	}

	return withoutCode(ctx, fmt.Errorf("api: unreachable"), 0)
}

func e404(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return withCode(withFail(ctx, estub("api: method %s", http.StatusNotFound)), http.StatusNotFound)
}

func e405(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return withCode(withFail(ctx, estub("api: %s", http.StatusMethodNotAllowed)), http.StatusMethodNotAllowed)
}

func estub(format string, code int) error {
	return fmt.Errorf(format, strings.ToLower(http.StatusText(code)))
}

func writeJSON(ctx context.Context, w http.ResponseWriter, code int, i interface{}) (int64, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return 0, err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Powered-By", runtime.Version())
	w.Header().Set("X-Request-ID", uuidFromContext(ctx))
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
	size := int64(n)

	_, err = w.Write([]byte("\n"))
	if err != nil {
		return 0, err
	}
	size++

	return size, nil
}
