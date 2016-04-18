package core

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"internal/compress/gzutil"

	"github.com/spkg/bom"
)

func stringOK() string {
	return http.StatusText(http.StatusOK)
}

func readClose(r io.ReadCloser) ([]byte, error) {
	defer func(c io.Closer) {
		_ = c.Close()
	}(r)

	return ioutil.ReadAll(r)
}

func isTypeGzip(b []byte) bool {
	return gzutil.IsGzipInString(http.DetectContentType(b))
}

func isTypeUTF8(b []byte) bool {
	return strings.Contains(http.DetectContentType(b), "text/plain; charset=utf-8")
}

func mendIfGzip(b []byte) ([]byte, error) {
	if isTypeGzip(b) {
		unz, err := gzutil.Gunzip(b)
		if err != nil {
			return nil, err
		}
		return unz, nil
	}

	return b, nil
}

func mendIfUTF8(b []byte) ([]byte, error) {
	if isTypeUTF8(b) {
		return bom.Clean(b), nil
	}

	return b, nil
}

func readMendClose(r io.ReadCloser) ([]byte, error) {
	b, err := readClose(r)
	if err != nil {
		return nil, err
	}

	b, err = mendIfGzip(b)
	if err != nil {
		return nil, err
	}

	b, err = mendIfUTF8(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func isEmpty(v []interface{}) bool {
	for i := range v {
		if v[i] != nil {
			return false
		}
	}
	return true
}
