package core

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"internal/util/gzutil"

	"github.com/spkg/bom"
)

func stringOK() string {
	return http.StatusText(http.StatusOK)
}

func btsToMD5(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

func strToMD5(s string) string {
	return btsToMD5([]byte(s))
}

func btsToSHA1(b []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(b))
}

func strToSHA1(s string) string {
	return btsToSHA1([]byte(s))
}

func readClose(r io.ReadCloser) ([]byte, error) {
	defer func() { _ = r.Close() }()
	return ioutil.ReadAll(r)
}

func isTypeGzip(b []byte) bool {
	return gzutil.IsGzipInString(http.DetectContentType(b))
}

func isTypeUTF8(b []byte) bool {
	return strings.Contains(http.DetectContentType(b), "text/plain; charset=utf-8")
}

func gunzip(b []byte) ([]byte, error) {
	b, err := gzutil.Gunzip(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func mendIfGzip(b []byte) ([]byte, error) {
	if isTypeGzip(b) {
		return gunzip(b)
	}

	return b, nil
}

func mendIfUTF8(b []byte) ([]byte, error) {
	if isTypeUTF8(b) {
		return bom.Clean(b), nil
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

func makeReadCloser(b []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(b))
}
