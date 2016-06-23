package api

import (
	"context"
	"net/http"
)

func pipeFail(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		err := failFromCtx(ctx)
		if err != nil {
			if code := codeFromCtx(ctx); code != 0 {
				var size int64
				size, err = writeResp(ctx, w, int(code), "err: "+err.Error())
				if err != nil {
					ctx = ctxWithFail(ctx, err)
				}
				ctx = ctxWithSize(ctx, size)
			}
		}
		h(ctx, w, r)
	}
}
