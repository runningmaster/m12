package api

import (
	"net/http"

	"internal/gzpool"
)

func pipeGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := failFromCtx(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if gzpool.IsGzipInString(r.Header.Get("Content-Encoding")) {
			z, err := gzpool.GetReader()
			if err != nil {
				ctx = ctxWithFail(ctx, err)
			}
			defer func() { _ = gzpool.PutReader(z) }()
			err = z.Reset(r.Body)
			if err != nil {
				ctx = ctxWithFail(ctx, err)
			}
			r.Body = z
		}

		if gzpool.IsGzipInString(r.Header.Get("Accept-Encoding")) {
			z, err := gzpool.GetWriter()
			if err != nil {
				ctx = ctxWithFail(ctx, err)
			}
			defer func() { _ = gzpool.PutWriter(z) }()
			z.Reset(w)
			w = gzpool.ResponseWriter{Writer: z, ResponseWriter: w}
			w.Header().Add("Vary", "Accept-Encoding")
			w.Header().Set("Content-Encoding", "gzip")
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
