package api

import "net/http"

func pipeFail(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
