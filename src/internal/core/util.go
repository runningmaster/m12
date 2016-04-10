package core

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"internal/compress/gzip"

	"github.com/spkg/bom"
)

func stringOK() string {
	return http.StatusText(http.StatusOK)
}

func readBody(r *http.Request) ([]byte, error) {
	defer func(c io.Closer) {
		_ = c.Close()
	}(r.Body)

	return ioutil.ReadAll(r.Body)
}

func isTypeGzip(b []byte) bool {
	return strings.Contains(http.DetectContentType(b), "gzip")
}

func isTypeUTF8(b []byte) bool {
	return strings.Contains(http.DetectContentType(b), "text/plain; charset=utf-8")
}

func mendIfGzip(b []byte) ([]byte, error) {
	if isTypeGzip(b) {
		unz, err := gzip.Gunzip(b)
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
