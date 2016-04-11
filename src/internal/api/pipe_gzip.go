package api

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"

	"internal/compress/gzutil"

	"github.com/klauspost/compress/gzip"
	"golang.org/x/net/context"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}

	n, err := w.Writer.Write(b)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (w gzipResponseWriter) Flush() error {
	return w.Writer.(*gzip.Writer).Flush()
}

func (w gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	c, rw, err := w.ResponseWriter.(http.Hijacker).Hijack()
	if err != nil {
		return c, rw, err
	}

	return c, rw, nil
}

func (w *gzipResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func gzipInContentEncoding(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
}
func gzipInAcceptEncoding(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func pipeGzip(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if gzipInContentEncoding(r) {
			z, err := gzutil.GetReader()
			if err != nil {
				// FIXME log err
			}
			defer gzutil.PutReader(z)
			if err := z.Reset(r.Body); err != nil {
				// FIXME log err
			}
			r.Body = z
		}

		if gzipInAcceptEncoding(r) {
			z, err := gzutil.GetWriter()
			if err != nil {
				// FIXME TODO log err
			}
			defer gzutil.PutWriter(z)
			z.Reset(w)
			w = gzipResponseWriter{Writer: z, ResponseWriter: w}
			w.Header().Add("Vary", "Accept-Encoding")
			w.Header().Set("Content-Encoding", "gzip")
		}

		h(ctx, w, r)
	}
}
