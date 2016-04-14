package core

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"internal/context/ctxutil"
	"internal/database/s3util"

	"golang.org/x/net/context"
)

/*
CheckGzip (if not then fail)
GetHeader [if not exist then find query param (Form)]
Calc body hash md5
sha1:	af064923bbf2301596aac4c273ba32178ebc4a96
md5:	b0804ec967f48520697662a204f5fe72
Put "H.Tag" "H.UUID"+.gz
SendToChannel
*/

// Upld sends data to s3 interface
func Upld(ctx context.Context, r *http.Request) (interface{}, error) {
	m, err := makeMeta(ctxutil.MetaFromContext(ctx))
	if err != nil {
		return nil, err
	}

	m.ID = ctxutil.IDFromContext(ctx)
	m.IP = ctxutil.IPFromContext(ctx)
	m.Key = ctxutil.AuthFromContext(ctx)
	m.SrcCE = r.Header.Get("Content-Encoding")
	m.SrcCT = r.Header.Get("Content-Type")

	p, err := packMeta(m)
	if err != nil {
		return nil, err
	}

	defer func(c io.Closer) {
		_ = c.Close()
	}(r.Body)

	if err := s3util.Put("stream-in", m.ID+".gz", r.Body, p); err != nil {
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
