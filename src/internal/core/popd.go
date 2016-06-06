package core

import (
	"net/http"
)

var Popd = &popd{}

type popd struct {
	meta []byte
}

func (p *popd) WriteHeader(h http.Header) {
	h.Set("Content-Encoding", "gzip")
	h.Set("Content-Meta", string(p.meta))
}

func (p *popd) Work(data []byte) (interface{}, error) {
	meta, data, err := popMetaData(data)
	if err != nil {
		return nil, err
	}
	p.meta = meta

	return data, nil
}
