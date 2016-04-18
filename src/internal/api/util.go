package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"

	"internal/context/ctxutil"

	"golang.org/x/net/context"
)

func with200(ctx context.Context, w http.ResponseWriter, res interface{}) context.Context {
	size, err := writeJSON(ctx, w, http.StatusOK, res)
	if err != nil {
		return withoutCode(ctx, err, size)
	}
	return ctxutil.WithCode(ctxutil.WithSize(ctx, size), http.StatusOK)
}

func with500(ctx context.Context, err error) context.Context {
	return ctxutil.WithCode(ctxutil.WithFail(ctx, err), http.StatusInternalServerError)
}

func withoutCode(ctx context.Context, err error, size int64) context.Context {
	return ctxutil.WithSize(ctxutil.WithFail(ctx, err), size)
}

func writeJSON(ctx context.Context, w http.ResponseWriter, code int, i interface{}) (int64, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return 0, err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Powered-By", runtime.Version())
	w.Header().Set("X-Request-ID", ctxutil.IDFromContext(ctx))
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
