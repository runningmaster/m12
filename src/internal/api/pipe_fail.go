package api

import (
	"net/http"

	"internal/errors"

	"golang.org/x/net/context"
)

func pipeFail(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if err := failFromContext(ctx); err != nil {
			if code := codeFromContext(ctx); code != 0 {
				var size int64
				if size, err = writeJSON(ctx, w, http.StatusInternalServerError, errors.Locus(err).Error()); err != nil {
					ctx = withFail(ctx, err)
				}
				ctx = withSize(ctx, size)
			}
		}
		h(ctx, w, r)
	}
}
