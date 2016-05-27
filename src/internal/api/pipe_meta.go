package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"internal/gzutil"

	"golang.org/x/net/context"
)

func pipeMeta(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var err error

		err = mustHeaderGzip(r.Header)
		if err != nil {
			goto fail
		}

		err = mustHeaderJSON(r.Header)
		if err != nil {
			goto fail
		}

		err = mustHeaderUTF8(r.Header)
		if err != nil {
			goto fail
		}

		err = mustHeaderMETA(r.Header)
		if err != nil {
			goto fail
		}

		// --------------------

		err = injectIntoMETA(r.Header)
		if err != nil {
			goto fail
		}

		h(ctx, w, r)
		return // success
	fail:
		h(ctxWithCode(ctxWithFail(ctx, err), http.StatusInternalServerError), w, r)
	}
}

func mustHeaderGzip(h http.Header) error {
	if !gzutil.IsGzipInString(h.Get("Content-Encoding")) {
		return fmt.Errorf("api: content-encoding must contain gzip")
	}
	return nil
}

func mustHeaderJSON(h http.Header) error {
	if !strings.Contains(h.Get("Content-Type"), "application/json") {
		return fmt.Errorf("api: content-type must contain application/json")
	}
	return nil
}

func mustHeaderUTF8(h http.Header) error {
	if !strings.Contains(h.Get("Content-Type"), "charset=utf-8") {
		return fmt.Errorf("api: content-type must contain charset=utf-8")
	}
	return nil
}

func mustHeaderMETA(h http.Header) error {
	if len(h.Get("Content-Meta")) == 0 {
		return fmt.Errorf("api: content-meta must contain value")
	}
	return nil
}

func injectIntoMETA(h http.Header) error {
	meta := h.Get("Content-Meta")
	b, err := base64.StdEncoding.DecodeString(s)

	//m.ID = ctxutil.IDFromContext(ctx)
	//m.IP = ctxutil.IPFromContext(ctx)
	//m.Auth = ctxutil.AuthFromContext(ctx)
	//m.Time = time.Now().Unix()

}
