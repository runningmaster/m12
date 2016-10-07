package pipe

import (
	"bytes"
	"io"
	"net/http"

	"internal/core/ctxt"
)

func Wrap(v interface{}) handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() { _ = r.Body.Close() }()

			ctx := r.Context()
			err := ctxt.FailFrom(ctx)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			var res interface{}
			var buf = new(bytes.Buffer)
			var n int64
			if r.Method == "POST" {
				n, err = io.Copy(buf, r.Body)
				if err != nil {
					goto exit
				}
			}
			ctx = ctxt.WithClen(ctx, n)

			switch h := v.(type) {
			case func(http.ResponseWriter, *http.Request):
				h(w, r) // stdh
				ctx = r.Context()
			case func([]byte, http.Header, http.Header) (interface{}, error):
				res, err = h(buf.Bytes(), r.Header, w.Header())
			case func([]byte) (interface{}, error):
				res, err = h(buf.Bytes())
			case func() (interface{}, error):
				res, err = h()
			default:
				panic("pipe: unknown handler")
			}

		exit:
			if err != nil {
				ctx = ctxt.WithFail(ctx, err)
			}
			if res != nil {
				ctx = ctxt.WithData(ctx, res)
			}
			*r = *r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
