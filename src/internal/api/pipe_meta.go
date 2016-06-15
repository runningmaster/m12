package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"internal/gzpool"

	"golang.org/x/net/context"
)

func pipeMeta(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		err := validateHeader(r.Header)
		if err != nil {
			goto fail
		}

		err = injectIntoMETA(ctx, r.Header)
		if err != nil {
			goto fail
		}

		h(ctx, w, r)
		return // success
	fail:
		h(ctxWithCode(ctxWithFail(ctx, err), http.StatusInternalServerError), w, r)
	}
}

func validateHeader(h http.Header) error {
	err := mustHeaderGzip(h)
	if err != nil {
		return err
	}

	err = mustHeaderJSON(h)
	if err != nil {
		return err
	}

	err = mustHeaderUTF8(h)
	if err != nil {
		return err
	}

	return mustHeaderMETA(h)

}

func mustHeaderGzip(h http.Header) error {
	if !gzpool.IsGzipInString(h.Get("Content-Encoding")) {
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

func injectIntoMETA(ctx context.Context, h http.Header) error {
	info := fmt.Sprintf(`{ "uuid": %q, "host": %q, "user": %q, "auth": %q, "time": %d, `,
		uuidFromCtx(ctx),
		hostFromCtx(ctx),
		userFromCtx(ctx),
		authFromCtx(ctx),
		timeFromCtx(ctx).Unix(),
	)

	meta, err := base64.StdEncoding.DecodeString(h.Get("Content-Meta"))
	if err != nil {
		return err
	}

	meta = bytes.Replace(bytes.TrimSpace(meta), []byte("{"), []byte(info), -1)
	h.Set("Content-Meta", string(meta))
	return nil
}
