package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

func pipeResp(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := failFromCtx(ctx)
		if err != nil {
			if code := codeFromCtx(ctx); code != 0 {
				var size int64
				size, err = writeResp(ctx, w, int(code), err.Error())
				if err != nil {
					ctx = ctxWithFail(ctx, err)
				}
				ctx = ctxWithSize(ctx, size)
			}
		}
		h(w, r.WithContext(ctx))
	}
}

func writeResp(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	var out []byte
	var err error
	if w.Header().Get("Content-Type") == "gzip" {
		var ok bool
		if out, ok = data.([]byte); !ok {
			return 0, fmt.Errorf("unknown data")
		}
	} else {
		if true { // FIXME (flag?)
			out, err = json.Marshal(data)
		} else {
			out, err = json.MarshalIndent(data, "", "\t")
		}
		if err != nil {
			return 0, err
		}
	}

	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Powered-By", runtime.Version())
	w.Header().Set("X-Request-ID", uuidFromCtx(ctx))
	w.WriteHeader(code)

	n, err := w.Write(out)
	if err != nil {
		ctx = ctxWithFail(ctx, err)
	}

	r = r.WithContext(ctxWithCode(ctxWithSize(ctx, n), code))
}

	n, err := writeResp(ctx, w, http.StatusOK, res)
	if err != nil {
		return ctxWithSize(ctxWithFail(ctx, err), n)
	}

	n, err := writeResp(ctx, w, http.StatusOK, s)
	if err != nil {
		r = r.WithContext(ctxWithSize(ctxWithFail(r.Context(), err), n))
		return
	}