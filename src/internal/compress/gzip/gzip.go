package gzip

import (
	"bytes"
	"io"
	"io/ioutil"
	"sync"

	"internal/errors"

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
				return errors.Locus(err)
			}
			return z
		},
	}

	writerPool = sync.Pool{
		New: func() interface{} {
			z, err := gzip.NewWriterLevel(ioutil.Discard, gzip.DefaultCompression)
			if err != nil {
				return errors.Locus(err)
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
	return nil, errors.Locusf("gzip: unreachable")
}

// CloseReader closes reader and puts it back to the pool.
func CloseReader(c io.Closer) {
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
	return nil, errors.Locusf("gzip: unreachable")
}

// CloseWriter closes writer and puts it back to the pool.
func CloseWriter(c io.Closer) {
	_ = c.Close()
	writerPool.Put(c)
}

// Gunzip decodes gzip-bytes.
func Gunzip(data []byte) ([]byte, error) {
	z, err := GetReader()
	if err != nil {
		return nil, errors.Locus(err)
	}
	defer CloseReader(z)

	if err = z.Reset(bytes.NewReader(data)); err != nil {
		return nil, errors.Locus(err)
	}

	var out []byte
	if out, err = ioutil.ReadAll(z); err != nil {
		return nil, errors.Locus(err)
	}

	return out, nil
}
