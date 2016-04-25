package core

import (
	"io"
	"net/http"

	"golang.org/x/net/context"
)

func Popd(_ context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	meta, data, err := mineMetaData(r.Body)
	if err != nil {
		return nil, err
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Meta", string(meta))

	return data, nil
}

func mineMetaData(rc io.ReadCloser) ([]byte, []byte, error) {
	b, err := readClose(rc)
	if err != nil {
		return nil, err
	}

	p := pathS3{}
	err = p.initFromJSON(b)
	if err != nil {
		return nil, err
	}

	o, err := popObject(p.Backet, p.Object)
	if err != nil {
		return nil, err
	}

	return untarMetaData(o)
}
