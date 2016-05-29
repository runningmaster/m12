package core

import (
	"net/http"

	"internal/s3"
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
	o, err := s3.PopObjectUnmarshal(data)
	if err != nil {
		return nil, err
	}

	meta, data, err := untarMetaData(o)
	if err != nil {
		return nil, err
	}
	p.meta = meta

	return data, nil
}
