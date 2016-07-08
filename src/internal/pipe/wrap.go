package pipe

import (
	"bytes"
	"io"
	"net/http"

	"internal/ctxutil"
)

func Wrap(v interface{}) handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() { _ = r.Body.Close() }()

			ctx := r.Context()
			err := ctxutil.FailFrom(ctx)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			var res interface{}
			switch h := v.(type) {
			case http.Handler:
				h.ServeHTTP(w, r)
				ctx = r.Context()
			case func(http.ResponseWriter, *http.Request):
				h(w, r)
				ctx = r.Context()
			case func([]byte, http.Header, http.Header) (interface{}, error):
				var buf = new(bytes.Buffer)
				n, err := io.Copy(buf, r.Body)
				if err != nil {
					ctx = ctxutil.WithFail(ctx, err)
				} else {
					res, err = h(buf.Bytes(), r.Header, w.Header())
					if err != nil {
						ctx = ctxutil.WithFail(ctx, err)
					}
				}
				ctx = ctxutil.WithClen(ctx, n)
			default:
				panic("pipe: unknown handler")
			}

			if res != nil {
				ctx = ctxutil.WithData(ctx, res)
			}

			*r = *r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
