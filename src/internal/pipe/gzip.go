package pipe

import (
	"net/http"

	"internal/ctxutil"
	"internal/gzip"
)

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := ctxutil.FailFrom(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if gzip.InString(r.Header.Get("Content-Encoding")) {
			z, err := gzip.GetReader()
			if err != nil {
				ctx = ctxutil.WithFail(ctx, err)
			}
			defer func() { _ = gzip.PutReader(z) }()
			err = z.Reset(r.Body)
			if err != nil {
				ctx = ctxutil.WithFail(ctx, err)
			}
			r.Body = z
		}

		if gzip.InString(r.Header.Get("Accept-Encoding")) {
			z, err := gzip.GetWriter()
			if err != nil {
				ctx = ctxutil.WithFail(ctx, err)
			}
			defer func() { _ = gzip.PutWriter(z) }()
			z.Reset(w)
			w = gzip.ResponseWriter{Writer: z, ResponseWriter: w}
			w.Header().Add("Vary", "Accept-Encoding")
			w.Header().Set("Content-Encoding", "gzip")
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
