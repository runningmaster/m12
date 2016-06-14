package core

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"net/http"
	"strings"

	"github.com/spkg/bom"
)

func btsToMD5(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

func btsToSHA1(b []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(b))
}

func strToSHA1(s string) string {
	return btsToSHA1([]byte(s))
}

func isUTF8(b []byte) bool {
	return strings.Contains(http.DetectContentType(b), "text/plain; charset=utf-8")
}

func mendIfUTF8(b []byte) ([]byte, error) {
	if isUTF8(b) {
		return bom.Clean(b), nil
	}

	return b, nil
}
