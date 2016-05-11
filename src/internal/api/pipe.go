package api

import (
	"net/http"

	"golang.org/x/net/context"
)

func use(pipes ...handlerPipe) http.Handler {
	h := func(context.Context, http.ResponseWriter, *http.Request) { /* dummy */ }
	for i := len(pipes) - 1; i >= 0; i-- {
		h = pipes[i](h) // note: nill will generate panic
	}
	return handlerFunc(h)
}

func pipe(h handlerFuncCtx) handlerPipe {
	return func(next handlerFunc) handlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			err := failFromContext(ctx)
			if err != nil {
				next(ctx, w, r)
				return
			}
			next(h(ctx, w, r), w, r)
		}
	}
}
