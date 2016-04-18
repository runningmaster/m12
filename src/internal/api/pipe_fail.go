package api

import (
	"net/http"

	"internal/context/ctxutil"

	"golang.org/x/net/context"
)

func pipeFail(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		err := ctxutil.FailFromContext(ctx)
		if err != nil {
			if code := ctxutil.CodeFromContext(ctx); code != 0 {
				var size int64
				size, err = writeJSON(ctx, w, http.StatusInternalServerError, err.Error())
				if err != nil {
					ctx = ctxutil.WithFail(ctx, err)
				}
				ctx = ctxutil.WithSize(ctx, size)
			}
		}
		h(ctx, w, r)
	}
}
