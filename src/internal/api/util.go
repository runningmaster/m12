package api

import (
	"bytes"
	"encoding/json"
	"net/http"

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

	w.Header().Set("X-Request-ID", ctxutil.IDFromContext(ctx))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if true { // FIXME (flag?)
		var tmp bytes.Buffer
		err = json.Indent(&tmp, b, "", "\t")
		if err != nil {
			return 0, err
		}
		b = tmp.Bytes()
	}

	var (
		n    int
		size int64
	)

	if n, err = w.Write(b); err != nil {
		return int64(n), err
	}
	size = size + int64(n)

	if _, err = w.Write([]byte("\n")); err != nil {
		return 0, err
	}
	size++

	return size, nil
}
