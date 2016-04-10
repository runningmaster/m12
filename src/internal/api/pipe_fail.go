package api

import (
	"net/http"

	"internal/context/ctxutil"

	"golang.org/x/net/context"
)

func pipeFail(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if err := ctxutil.FailFromContext(ctx); err != nil {
			if code := ctxutil.CodeFromContext(ctx); code != 0 {
				var size int64
				if size, err = writeJSON(ctx, w, http.StatusInternalServerError, err.Error()); err != nil {
					ctx = ctxutil.WithFail(ctx, err)
				}
				ctx = ctxutil.WithSize(ctx, size)
			}
		}
		h(ctx, w, r)
	}
}
