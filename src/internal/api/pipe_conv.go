package api

import (
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"golang.org/x/net/context"
)

func pipeConv(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		convHost(r)
		convAuth(r)
		convHTag(r)
		h(ctx, w, r)
	}
}

func convHost(r *http.Request) {
	v := r.Header.Get("X-Morion-Skynet-Origin-IP")
	if v != "" {
		r.Header.Set("X-Forwarded-For", v)
		r.Header.Set("X-Real-IP", v)
	}
}

func convAuth(r *http.Request) {
	// V2 X-Morion-Skynet-Key: value
	v := r.Header.Get("X-Morion-Skynet-Key")
	if v == "" {
		// V1 ?key=value
		v = r.URL.Query().Get("key")
	}
	if v != "" {
		r.SetBasicAuth("api", fmt.Sprintf("key-%s", v))
	}
}

func convHTag(r *http.Request) {
	// V2 X-Morion-Skynet-Tag: value
	v := r.Header.Get("X-Morion-Skynet-Tag")
	if v == "" {
		// V1 Content-Type: application/json; charset=utf-8; hashtag=value
		v = r.Header.Get("Content-Type")
		_, p, err := mime.ParseMediaType(v)
		if err != nil {
			v = ""
		} else {
			v = p["hashtag"]
		}
	}

	if v != "" { // FIXME: remove test (!)
		v = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`{ "test": true, "ctag": "conv", "htag": %q }`, strings.ToLower(v))))
		r.Header.Set("Content-Meta", v)
		r.Header.Set("Content-Type", "application/json; charset=utf-8")
		r.Header.Set("Content-Encoding", "gzip")
	}
}
