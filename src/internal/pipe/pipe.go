package pipe

import (
	"internal/ctxutil"
	"net/http"
)

type handler func(h http.Handler) http.Handler

func Use(pipes ...handler) http.Handler {
	var h http.Handler
	for i := len(pipes) - 1; i >= 0; i-- {
		h = pipes[i](h) // note: nill will generate panic
	}
	return h
}

func Pipe(v interface{}) handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := ctxutil.FailFrom(r.Context())
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			switch h := v.(type) {
			case http.Handler:
				h.ServeHTTP(w, r)
			case func(http.ResponseWriter, *http.Request):
				h(w, r)
			default:
				panic("api: unknown handler")
			}

			next.ServeHTTP(w, r)
		})
	}
}
