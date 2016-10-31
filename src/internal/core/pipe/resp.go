package pipe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"internal/core/pref"
)

func Resp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		data, err := dataFrom(ctx), failFrom(ctx)
		if err != nil {
			data = err.Error()
		}

		// workaround for stdh
		if stdhFrom(ctx) {
			next.ServeHTTP(w, r)
			return
		}

		n, err := writeResp(w, uuidFrom(ctx), codeFrom(ctx), data)
		if err != nil {
			ctx = withFail(ctx, err)
		}

		ctx = withSize(ctx, int64(n))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeResp(w http.ResponseWriter, uuid string, code int, data interface{}) (int, error) {
	var out []byte
	var err error
	// FIXME
	if data != nil {
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
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Powered-By", fmt.Sprintf("go version %s", runtime.Version()))
	w.Header().Set("X-Request-ID", uuid)
	w.WriteHeader(code)

	if len(out) > 0 {
		return w.Write(out)
	}

	return 0, nil

	//	_, _ = w.Write([]byte("\n")) // (?)
	//	return n + 1, nil
}
