package pipe

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"internal/compress/gziputil"
)

func Meta(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := failFrom(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		err = injectMeta(ctx, r.Header)
		if err != nil {
			ctx = withFail(ctx, err)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
		Time string `json:"time,omitempty"`
		Unix int64  `json:"unix,omitempty"`
	}{}

	v.UUID = uuidFrom(ctx)
	v.Auth.ID = authFrom(ctx)
	v.Host = hostFrom(ctx)
	v.User = userFrom(ctx)
	v.Time = timeFrom(ctx).String()
	v.Unix = timeFrom(ctx).Unix()

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
	if !gziputil.InString(h.Get("Content-Encoding")) {
		return fmt.Errorf("pipe: content-encoding must contain gzip")
	}
	return nil
}

func mustHeaderJSON(h http.Header) error {
	if !strings.Contains(h.Get("Content-Type"), "application/json") {
		return fmt.Errorf("pipe: content-type must contain application/json")
	}
	return nil
}

func mustHeaderUTF8(h http.Header) error {
	if !strings.Contains(h.Get("Content-Type"), "charset=utf-8") {
		return fmt.Errorf("pipe: content-type must contain charset=utf-8")
	}
	return nil
}

func mustHeaderMETA(h http.Header) error {
	if h.Get("Content-Meta") == "" {
		return fmt.Errorf("pipe: content-meta must contain value")
	}
	return nil
}
