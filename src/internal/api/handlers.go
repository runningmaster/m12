package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"internal/errors"
	"internal/flag"
	"internal/version"

	"golang.org/x/net/context"
)

type (
	handlerFunc    func(context.Context, http.ResponseWriter, *http.Request)
	handlerFuncCtx func(context.Context, http.ResponseWriter, *http.Request) context.Context
	handlerPipe    func(h handlerFunc) handlerFunc
	coreFunc       func([]byte) (interface{}, error)

	bundle struct {
		h http.Handler
		f coreFunc
	}
)

func (f handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(context.Background(), w, r)
}

func root(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	res := fmt.Sprintf("%s %s", version.Stamp.AppName(), version.Stamp.Extended())
	return with200(ctx, w, res)
}

func exec(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	var (
		b   []byte
		err error
	)

	if r.Method == "POST" {
		defer func(c io.Closer) {
			_ = c.Close()
		}(r.Body)

		if b, err = ioutil.ReadAll(r.Body); err != nil {
			return with500(ctx, errors.Locus(err))
		}
	}

	var res interface{}
	if f, ok := mapFuncs[r.URL.Path]; ok {
		res, err = f(b)
		if err != nil {
			return with500(ctx, errors.Locus(err))
		}
	} else {
		panic(errors.Locusf("exec: unreachable"))
	}

	return with200(ctx, w, res)
}

func stdh(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	if !flag.Debug {
		return with500(ctx, errors.Locusf("flag debug not found"))
	}

	if h, p := http.DefaultServeMux.Handler(r); p != "" {
		h.ServeHTTP(w, r)
		return withCode(withSize(ctx, 0), http.StatusOK) // TODO: wrap w to get real size
	}

	return withoutCode(ctx, errors.Locusf("debug: unreachable"), 0)
}

func e404(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return withCode(withFail(ctx, estub("api: method %s", http.StatusNotFound)), http.StatusNotFound)
}

func e405(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return withCode(withFail(ctx, estub("api: %s", http.StatusMethodNotAllowed)), http.StatusMethodNotAllowed)
}

func estub(format string, code int) error {
	return errors.Locusf(format, strings.ToLower(http.StatusText(code)))
}
