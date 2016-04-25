package core

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
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

func mendIfGzipUTF8(b []byte) ([]byte, error) {
	u, err := mendIfGzip(b)
	if err != nil {
		return nil, err
	}

	return mendIfUTF8(u)
}

func readMendClose(r io.ReadCloser) ([]byte, error) {
	b, err := readClose(r)
	if err != nil {
		return nil, err
	}

	return mendIfGzipUTF8(b)
}

func isEmpty(v []interface{}) bool {
	for i := range v {
		if v[i] != nil {
			return false
		}
	}
	return true
}
