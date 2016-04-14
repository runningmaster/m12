package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"internal/context/ctxutil"

	"golang.org/x/net/context"
)

func pipeMeta(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		meta, err := getMeta(r)
		if err != nil {
			goto fail
		}
		h(ctxutil.WithMeta(ctx, meta), w, r)
		return
	fail:
		h(ctxutil.WithCode(ctxutil.WithFail(ctx, err), http.StatusInternalServerError), w, r)
	}
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
