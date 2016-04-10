package core

import (
	"io"
	"net/http"

	"internal/context/ctxutil"
	"internal/net/s3"

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
	//	var (
	//		b   []byte
	//		err error
	//	)
	//	if b, err = readBody(r); err != nil {
	//		return nil, err
	//	}
	//
	//	if !isTypeGzip(b) {
	//		return nil, fmt.Errorf("core: s3: gzip not found")
	//	}

	_ = s3.MkB("input")
	//	if err != nil {
	//		return nil, err
	//	}
	id := ctxutil.IDFromContext(ctx)

	if err := s3.Put("input", id+".gz", r.Body /*bytes.NewBuffer(b)*/, "{"+id+"}"); err != nil {
		return nil, err
	}
	defer func(c io.Closer) {
		_ = c.Close()
	}(r.Body)
	// s3.Get and check content type

	return "OK: " + id, nil
}
