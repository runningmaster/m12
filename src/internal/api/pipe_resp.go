package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"internal/pref"
)

func pipeResp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		data, err := dataFromCtx(ctx), failFromCtx(ctx)
		if err != nil {
			data = err.Error()
		}

		// workaround for stdh
		if sizeFromCtx(ctx) != 0 {
			next.ServeHTTP(w, r)
			return
		}

		n, err := writeResp(w, uuidFromCtx(ctx), codeFromCtx(ctx), data)
		if err != nil {
			ctx = ctxWithFail(ctx, err)
		}

		ctx = ctxWithSize(ctx, int64(n))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeResp(w http.ResponseWriter, uuid string, code int, data interface{}) (int, error) {
	var out []byte
	var err error
	// FIXME
	if w.Header().Get("Content-Type") == "gzip" {
		var ok bool
		if out, ok = data.([]byte); !ok {
			return 0, fmt.Errorf("unknown data")
		}
	} else {
		if !pref.Debug { // FIXME (flag?)
			out, err = json.Marshal(data)
		} else {
			out, err = json.MarshalIndent(data, "", "\t")
		}
		if err != nil {
			return 0, err
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Powered-By", fmt.Sprintf("go version %s", runtime.Version()))
	w.Header().Set("X-Request-ID", uuid)
	w.WriteHeader(code)

	n, err := w.Write(out)
	if err != nil {
		return 0, err
	}

	_, _ = w.Write([]byte("\n")) // (?)
	return n + 1, nil
}
