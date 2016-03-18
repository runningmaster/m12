package api

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"

	"internal/errors"

	"github.com/klauspost/compress/gzip"
	"golang.org/x/net/context"
)

var (
	// gzipSimple is gzip empty string for init reader
	gzipSimple = []byte{
		0x1f, 0x8b, 0x8, 0x0, 0x0, 0x9, 0x6e, 0x88, 0x0, 0xff, 0x1,
		0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	}

	gzipReaderPool = sync.Pool{
		New: func() interface{} {
			z, err := gzip.NewReader(bytes.NewReader(gzipSimple))
			if err != nil {
				return errors.Locus(err)
			}
			return z
		},
	}

	gzipWriterPool = sync.Pool{
		New: func() interface{} {
			z, err := gzip.NewWriterLevel(ioutil.Discard, gzip.DefaultCompression)
			if err != nil {
				return errors.Locus(err)
			}
			return z
		},
	}
)

type gzipWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipWriter) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}

	n, err := w.Writer.Write(b)
	if err != nil {
		return n, errors.Locus(err)
	}

	return n, nil
}

func (w gzipWriter) Flush() error {
	return errors.Locus(w.Writer.(*gzip.Writer).Flush())
}

func (w gzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	c, rw, err := w.ResponseWriter.(http.Hijacker).Hijack()
	if err != nil {
		return c, rw, errors.Locus(err)
	}

	return c, rw, nil
}

func (w *gzipWriter) CloseNotify() <-chan bool {
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
			switch z := gzipReaderPool.Get().(type) {
			case *gzip.Reader:
				defer func() {
					_ = z.Close()
					gzipReaderPool.Put(z)
				}()
				if err := z.Reset(r.Body); err != nil {
					// FIXME TODO log err
				}
				r.Body = z
			case error:
				panic(z) // FIXME TODO
			}
		}

		if gzipInAcceptEncoding(r) {
			switch z := gzipWriterPool.Get().(type) {
			case *gzip.Writer:
				defer func() {
					_ = z.Close()
					gzipWriterPool.Put(z)
				}()
				z.Reset(w)
				w = gzipWriter{Writer: z, ResponseWriter: w}
				w.Header().Add("Vary", "Accept-Encoding")
				w.Header().Set("Content-Encoding", "gzip")
			case error:
				panic(z) // FIXME TODO
			}
		}

		h(ctx, w, r)
	}
}

func gunzip(data []byte) ([]byte, error) {
	switch z := gzipReaderPool.Get().(type) {
	case *gzip.Reader:
		defer func() {
			_ = z.Close()
			gzipReaderPool.Put(z)
		}()

		if err := z.Reset(bytes.NewReader(data)); err != nil {
			return nil, errors.Locus(err)
		}

		out, err := ioutil.ReadAll(z)
		if err != nil {
			return nil, errors.Locus(err)
		}

		return out, nil
	case error:
		return nil, z
	}
	return nil, errors.Locusf("gunzip: unreachable")
}
