package pipe

import (
	"net/http"

	"internal/ctxutil"
)

type handler func(h http.Handler) http.Handler

func Use(pipes ...handler) http.Handler {
	var h http.Handler
	for i := len(pipes) - 1; i >= 0; i-- {
		h = pipes[i](h) // note: nill will generate panic
	}
	return h
}

func Work(v interface{}) handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			err := ctxutil.FailFrom(ctx)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			switch h := v.(type) {
			case http.Handler:
				h.ServeHTTP(w, r)
			case func(http.ResponseWriter, *http.Request):
				h(w, r)
			case func() (interface{}, error):
				res, err := h()
				if err != nil {
					*r = *r.WithContext(ctxutil.WithFail(ctx, err))
					next.ServeHTTP(w, r)
					return
				}
				ctx = ctxutil.WithData(ctx, res)
				*r = *r.WithContext(ctx)
			default:
				panic("pipe: unknown handler")
			}

			next.ServeHTTP(w, r)
		})
	}
}
