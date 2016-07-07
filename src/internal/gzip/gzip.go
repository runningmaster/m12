package gzip

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/klauspost/compress/gzip"
)

var (
	// simple is gzip empty string for init reader
	simple = []byte{
		0x1f, 0x8b, 0x8, 0x0, 0x0, 0x9, 0x6e, 0x88, 0x0, 0xff, 0x1,
		0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	}

	readerPool = sync.Pool{
		New: func() interface{} {
			z, err := gzip.NewReader(bytes.NewReader(simple))
			if err != nil {
				return err
			}
			return z
		},
	}

	writerPool = sync.Pool{
		New: func() interface{} {
			z, err := gzip.NewWriterLevel(ioutil.Discard, gzip.DefaultCompression)
			if err != nil {
				return err
			}
			return z
		},
	}
)

// GetReader gets reader from pool.
func GetReader() (*gzip.Reader, error) {
	switch r := readerPool.Get().(type) {
	case *gzip.Reader:
		return r, nil
	case error:
		return nil, r
	}
	return nil, fmt.Errorf("gzip: unreachable")
}

// PutReader closes reader and puts it back to the pool.
func PutReader(c io.Closer) error {
	if c == nil {
		return nil
	}

	err := c.Close()
	if err != nil {
		return err
	}
	readerPool.Put(c)
	return nil
}

// GetWriter gets writer from pool.
func GetWriter() (*gzip.Writer, error) {
	switch w := writerPool.Get().(type) {
	case *gzip.Writer:
		return w, nil
	case error:
		return nil, w
	}
	return nil, fmt.Errorf("gzip: unreachable")
}

// PutWriter closes writer and puts it back to the pool.
func PutWriter(c io.Closer) error {
	if c == nil {
		return nil
	}

	err := c.Close()
	if err != nil {
		return err
	}
	writerPool.Put(c)
	return nil
}

// Compress compresses bytes.
func Compress(data []byte) ([]byte, error) {
	w, err := GetWriter()
	if err != nil {
		return nil, err
	}
	defer func() { _ = PutWriter(w) }()

	buf := new(bytes.Buffer)
	w.Reset(buf)

	_, err = io.Copy(w, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// MustGzip encodes bytes if not.
func MustCompress(data []byte) ([]byte, error) {
	if InString(http.DetectContentType(data)) {
		return data, nil
	}

	return Compress(data)
}

// Uncompress decompresses bytes.
func Uncompress(data []byte) ([]byte, error) {
	r, err := GetReader()
	if err != nil {
		return nil, err
	}
	defer func() { _ = PutReader(r) }()

	err = r.Reset(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Copy reads and ungzips from Reader to Writer
func Copy(dst io.Writer, src io.Reader) error {
	r, err := GetReader()
	if err != nil {
		return err
	}
	defer func() { _ = PutReader(r) }()

	err = r.Reset(src)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, r)
	return err
}

// InString returns true if gzip is mentioned in string
func InString(s string) bool {
	return strings.Contains(s, "gzip")
}

// ResponseWriter is gzip-wrapper for http.ResponseWriter
type ResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}

	return w.Writer.Write(b)
}

func (w ResponseWriter) Flush() error {
	return w.Writer.(*gzip.Writer).Flush()
}

func (w ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *ResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
