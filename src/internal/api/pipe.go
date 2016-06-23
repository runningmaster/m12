package api

import "net/http"

func use(pipes ...handlerPipe) http.Handler {
	h := func(http.ResponseWriter, *http.Request) { /* dummy */ }
	for i := len(pipes) - 1; i >= 0; i-- {
		h = pipes[i](h) // note: nill will generate panic
	}
	return http.HandlerFunc(h)
}

func pipe(h http.HandlerFunc) handlerPipe {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			err := failFromCtx(ctx)
			if err != nil {
				next(ctx, w, r)
				return
			}
			next(h(ctx, w, r), w, r)
		}
	}
}
