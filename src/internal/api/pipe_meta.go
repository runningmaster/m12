package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"internal/compress/gzutil"
	"internal/context/ctxutil"

	"golang.org/x/net/context"
)

func pipeMeta(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			m   string
			err error
		)
		if m, err = getMeta(r); err != nil {
			goto fail
		}

		if err = mustHeaderGzip(r); err != nil {
			goto fail
		}

		if err = mustHeaderJSON(r); err != nil {
			goto fail
		}

		if err = mustHeaderUTF8(r); err != nil {
			goto fail
		}

		h(ctxutil.WithMeta(ctx, m), w, r)
		return // success
	fail:
		h(ctxutil.WithCode(ctxutil.WithFail(ctx, err), http.StatusInternalServerError), w, r)
	}
}

func mustHeaderGzip(r *http.Request) error {
	if !gzutil.IsGzipInString(r.Header.Get("Content-Encoding")) {
		return fmt.Errorf("api: content-encoding must contain 'gzip'")
	}

	return nil
}

func mustHeaderJSON(r *http.Request) error {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("api: content-type must contain 'application/json'")
	}

	return nil
}

func mustHeaderUTF8(r *http.Request) error {
	if !strings.Contains(r.Header.Get("Content-Type"), "charset=utf-8") {
		return fmt.Errorf("api: content-type must contain 'charset=utf-8'")
	}

	return nil
}

func getMeta(r *http.Request) (string, error) {
	m := r.Header.Get("Content-Meta")
	if m == "" {
		return "", fmt.Errorf("api: content-meta not found")
	}

	b, err := base64.StdEncoding.DecodeString(m)
	if err != nil {
		return "", err
	}

	// check for correct json
	var v struct{}
	if err = json.Unmarshal(b, &v); err != nil {
		return "", fmt.Errorf("api: content-meta must be correct json: %s", err)
	}

	return string(b), nil
}
