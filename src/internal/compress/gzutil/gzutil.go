package gzutil

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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
func PutReader(c io.Closer) {
	_ = c.Close()
	readerPool.Put(c)
}

// GetWriter gets writer from pool.
func GetWriter() (*gzip.Writer, error) {
	switch w := readerPool.Get().(type) {
	case *gzip.Writer:
		return w, nil
	case error:
		return nil, w
	}
	return nil, fmt.Errorf("gzip: unreachable")
}

// PutWriter closes writer and puts it back to the pool.
func PutWriter(c io.Closer) {
	_ = c.Close()
	writerPool.Put(c)
}

// Gunzip decodes gzip-bytes.
func Gunzip(data []byte) ([]byte, error) {
	z, err := GetReader()
	if err != nil {
		return nil, err
	}
	defer PutReader(z)

	if err = z.Reset(bytes.NewReader(data)); err != nil {
		return nil, err
	}

	var out []byte
	if out, err = ioutil.ReadAll(z); err != nil {
		return nil, err
	}

	return out, nil
}

// NeedGzip returns true if gzip is mentioned in string
func NeedGzip(s string) bool {
	return strings.Contains(s, "gzip")
}
