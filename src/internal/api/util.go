package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"internal/errors"

	"github.com/spkg/bom"
)

func writeJSON(w http.ResponseWriter, code int, i interface{}) (int64, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return 0, errors.Locus(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if true { // FIXME (flag?)
		var tmp bytes.Buffer
		err = json.Indent(&tmp, b, "", "\t")
		if err != nil {
			return 0, errors.Locus(err)
		}
		b = tmp.Bytes()
	}

	var (
		n    int
		size int64
	)

	if n, err = w.Write(b); err != nil {
		return int64(n), errors.Locus(err)
	}
	size = size + int64(n)

	if _, err = w.Write([]byte("\n")); err != nil {
		return 0, errors.Locus(err)
	}
	size++

	return size, nil
}

func readAllAndClose(r io.ReadCloser) ([]byte, error) {
	defer func(c io.Closer) {
		_ = c.Close()
	}(r)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Locus(err)
	}

	return mendGzipAndBOM(b)
}

func mendGzipAndBOM(b []byte) ([]byte, error) {
	ct := http.DetectContentType(b)

	// Gzip
	if strings.Contains(ct, "gzip") {
		unz, err := gunzip(b)
		if err != nil {
			return nil, errors.Locus(err)
		}
		return unz, nil
	}

	// UTF BOM
	if strings.Contains(ct, "text/plain; charset=utf-8") {
		return bom.Clean(b), nil
	}

	return b, nil
}
