package api

import "net/http"

func use(pipes ...handlerPipe) http.Handler {
	var h http.Handler
	for i := len(pipes) - 1; i >= 0; i-- {
		h = pipes[i](h) // note: nill will generate panic
	}
	return h
}

func pipe(v interface{}) handlerPipe {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := failFromCtx(r.Context())
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
