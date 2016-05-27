package core

import (
	"net/http"

	"internal/s3"

	"golang.org/x/net/context"
)

var Popd = &popd{}

type popd struct {
	test string
}

func (u *popd) ReadHeader(h http.Header) {
	_ = h.Get("Content-Meta")
	u.test = h.Get("Content-Test")
}

func (u *popd) WriteHeader(h http.Header) {
	h.Set("Content-Test", u.test+" DEBUG")
}

func (u *popd) Work([]byte) (interface{}, error) {
	return nil, nil
}

func Popd2(_ context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	b, err := readClose(r.Body)
	if err != nil {
		return nil, err
	}

	o, err := s3.PopObjectByPathJSON(b)
	if err != nil {
		return nil, err
	}

	meta, data, err := untarMetaData(o)
	if err != nil {
		return nil, err
	}

	m, err := makeMetaFromJSON(meta)
	if err != nil {
		return nil, err
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Meta", m.packToBase64String()) // FIXME

	return data, nil
}
