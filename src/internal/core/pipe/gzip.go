package pipe

import (
	"net/http"

	"internal/compress/gziputil"
)

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := failFrom(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if gziputil.InString(r.Header.Get("Content-Encoding")) {
			z, err := gziputil.GetReader()
			if err != nil {
				ctx = withFail(ctx, err)
			}
			defer func() { _ = gziputil.PutReader(z) }()
			err = z.Reset(r.Body)
			if err != nil {
				ctx = withFail(ctx, err)
			}
			r.Body = z
		}

		if gziputil.InString(r.Header.Get("Accept-Encoding")) {
			z, err := gziputil.GetWriter()
			if err != nil {
				ctx = withFail(ctx, err)
			}
			defer func() { _ = gziputil.PutWriter(z) }()
			z.Reset(w)
			w = gziputil.ResponseWriter{Writer: z, ResponseWriter: w}
			w.Header().Add("Vary", "Accept-Encoding")
			w.Header().Set("Content-Encoding", "gzip")
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
