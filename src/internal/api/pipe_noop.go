package api

import (
	"net/http"

	"golang.org/x/net/context"
)

func pipeNoop(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		h(ctx, w, r)
	}
}
