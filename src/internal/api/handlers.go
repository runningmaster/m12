package api

import (
	"fmt"
	"net/http"
	"strings"

	"internal/context/ctxutil"
	"internal/core"
	"internal/flag"
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
	if !flag.Debug {
		return with500(ctx, fmt.Errorf("api: flag debug not found"))
	}

	if h, p := http.DefaultServeMux.Handler(r); p != "" {
		h.ServeHTTP(w, r)
		return ctxutil.WithCode(ctxutil.WithSize(ctx, 0), http.StatusOK) // TODO: wrap w to get real size
	}

	return withoutCode(ctx, fmt.Errorf("api: unreachable"), 0)
}

func e404(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return ctxutil.WithCode(ctxutil.WithFail(ctx, estub("api: method %s", http.StatusNotFound)), http.StatusNotFound)
}

func e405(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return ctxutil.WithCode(ctxutil.WithFail(ctx, estub("api: %s", http.StatusMethodNotAllowed)), http.StatusMethodNotAllowed)
}

func estub(format string, code int) error {
	return fmt.Errorf(format, strings.ToLower(http.StatusText(code)))
}
