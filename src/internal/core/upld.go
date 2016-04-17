package core

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"internal/context/ctxutil"

	"golang.org/x/net/context"
)

// Upld sends data to s3 interface
func Upld(ctx context.Context, r *http.Request) (interface{}, error) {
	m, err := makeMeta(ctxutil.MetaFromContext(ctx))
	if err != nil {
		return nil, err
	}

	m.ID = ctxutil.IDFromContext(ctx)
	m.IP = ctxutil.IPFromContext(ctx)
	m.Auth = ctxutil.AuthFromContext(ctx)
	m.HTag = strings.ToLower(m.HTag)
	m.SrcCE = r.Header.Get("Content-Encoding")
	m.SrcCT = r.Header.Get("Content-Type")

	p, err := packMeta(m)
	if err != nil {
		return nil, err
	}

	defer func(c io.Closer) {
		_ = c.Close()
	}(r.Body)

	if _, err = s3cli.PutObject(backetStreamIn, m.ID, r.Body, p); err != nil {
		return nil, err
	}

	// send ?

	return m.ID, nil
}

func makeMeta(s string) (meta, error) {
	var v meta
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return meta{}, err
	}
	return v, nil
}

func packMeta(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
