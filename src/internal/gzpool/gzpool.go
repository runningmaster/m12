package gzpool

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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
	err := c.Close()
	if err != nil {
		return err
	}
	writerPool.Put(c)
	return nil
}

// Gzip encodes bytes to gzip-bytes.
func Gzip(data []byte) ([]byte, error) {
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

	return buf.Bytes(), nil
}

// MustGzip encodes bytes to gzip-bytes if not.
func MustGzip(data []byte) ([]byte, error) {
	if IsGzipInString(http.DetectContentType(data)) {
		return data, nil
	}

	return Gzip(data)
}

// Gunzip decodes bytes from gzip-bytes.
func Gunzip(data []byte) ([]byte, error) {
	r, err := GetReader()
	if err != nil {
		return nil, err
	}
	defer func() { _ = PutReader(r) }()

	err = r.Reset(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		//BUG IS HERE
		return nil, err
	}

	return b, nil
}

// IsGzipInString returns true if gzip is mentioned in string
func IsGzipInString(s string) bool {
	return strings.Contains(s, "gzip")
}
