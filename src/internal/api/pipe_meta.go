package api

import (
	"fmt"
	"net/http"
	"strings"

	"internal/compress/gzutil"
	"internal/context/ctxutil"

	"golang.org/x/net/context"
)

func pipeMeta(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		err := mustHeaderMeta(r)
		if err != nil {
			goto fail
		}

		err = mustHeaderGzip(r)
		if err != nil {
			goto fail
		}

		err = mustHeaderJSON(r)
		if err != nil {
			goto fail
		}

		err = mustHeaderUTF8(r)
		if err != nil {
			goto fail
		}

		h(ctxutil.WithMeta(ctx, r.Header.Get("Content-Meta")), w, r)
		return // success
	fail:
		h(ctxutil.WithCode(ctxutil.WithFail(ctx, err), http.StatusInternalServerError), w, r)
	}
}

func mustHeaderGzip(r *http.Request) error {
	if !gzutil.IsGzipInString(r.Header.Get("Content-Encoding")) {
		return fmt.Errorf("api: content-encoding must contain gzip")
	}

	return nil
}

func mustHeaderJSON(r *http.Request) error {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("api: content-type must contain application/json")
	}

	return nil
}

func mustHeaderUTF8(r *http.Request) error {
	if !strings.Contains(r.Header.Get("Content-Type"), "charset=utf-8") {
		return fmt.Errorf("api: content-type must contain charset=utf-8")
	}

	return nil
}

func mustHeaderMeta(r *http.Request) error {
	if r.Header.Get("Content-Meta") == "" {
		return fmt.Errorf("api: header must contain content-meta")
	}

	return nil
}
