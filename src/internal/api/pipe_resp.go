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
		data, err := dataFromCtx(ctx), failFromCtx(ctx)
		if err != nil {
			data = err
		}

		uuid := uuidFromCtx(ctx)
		code := codeFromCtx(ctx)

		n, err := writeResp(w, uuid, code, data)
		if err != nil {
			ctx = ctxWithFail(ctx, err)
		}

		ctx = ctxWithSize(ctx, int64(n))
		h(w, r.WithContext(ctx))
	}
}

func writeResp(w http.ResponseWriter, uuid string, code int, data interface{}) (int, error) {
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Powered-By", runtime.Version())
	w.Header().Set("X-Request-ID", uuid)
	w.WriteHeader(code)

	return w.Write(out)
}
