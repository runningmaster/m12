package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

func writeJSON(ctx context.Context, w http.ResponseWriter, code int, i interface{}) (int64, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return 0, err
	}

	w.Header().Set("X-UUID", uuidFromContext(ctx))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if true { // FIXME (flag?)
		var tmp bytes.Buffer
		err = json.Indent(&tmp, b, "", "\t")
		if err != nil {
			return 0, err
		}
		b = tmp.Bytes()
	}

	var (
		n    int
		size int64
	)

	if n, err = w.Write(b); err != nil {
		return int64(n), err
	}
	size = size + int64(n)

	if _, err = w.Write([]byte("\n")); err != nil {
		return 0, err
	}
	size++

	return size, nil
}
