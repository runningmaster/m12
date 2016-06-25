package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"internal/gzpool"
)

func pipeMeta(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := failFromCtx(ctx)
		if err != nil {
			h(w, r)
			return
		}

		err = injectMeta(ctx, r.Header)
		if err != nil {
			ctx = ctxWithFail(ctx, err)
		}

		h(w, r.WithContext(ctx))
	}
}

func injectMeta(ctx context.Context, h http.Header) error {
	err := checkHeader(h)
	if err != nil {
		return err
	}

	meta, err := base64.StdEncoding.DecodeString(h.Get("Content-Meta"))
	if err != nil {
		return err
	}

	v := struct {
		UUID string `json:"uuid,omitempty"`
		Auth struct {
			ID string `json:"id,omitempty"`
		} `json:"auth,omitempty"`
		Host string `json:"host,omitempty"`
		User string `json:"user,omitempty"`
		Time int64  `json:"time,omitempty"`
	}{}

	v.UUID = uuidFromCtx(ctx)
	v.Auth.ID = authFromCtx(ctx)
	v.Host = hostFromCtx(ctx)
	v.User = userFromCtx(ctx)
	v.Time = timeFromCtx(ctx).Unix()

	m, err := json.Marshal(v)
	if err != nil {
		return err
	}

	meta = bytes.Replace(meta, []byte("{"), []byte(","), -1)
	meta = append(m[:len(m)-1], meta...)

	h.Set("Content-Meta", string(meta))
	return nil
}

func checkHeader(h http.Header) error {
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
