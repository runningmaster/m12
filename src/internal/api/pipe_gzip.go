package api

import (
	"bufio"
	"io"
	"net"
	"net/http"

	"internal/gzpool"

	"github.com/klauspost/compress/gzip"
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
			w = gzipResponseWriter{Writer: z, ResponseWriter: w}
			w.Header().Add("Vary", "Accept-Encoding")
			w.Header().Set("Content-Encoding", "gzip")
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}

	return w.Writer.Write(b)
}

func (w gzipResponseWriter) Flush() error {
	return w.Writer.(*gzip.Writer).Flush()
}

func (w gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *gzipResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
