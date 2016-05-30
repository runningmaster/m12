package core

import (
	"bytes"
	"net/http"

	"internal/s3"
)

var Putd = &putd{}

type putd struct {
	meta []byte
}

func (p *putd) ReadHeader(h http.Header) {
	p.meta = []byte(h.Get("Content-Meta"))
}

func (p *putd) Work(data []byte) (interface{}, error) {
	pos := bytes.IndexByte(p.meta, '.')
	uuid := string(p.meta[:pos])
	meta := p.meta[pos+1:]

	t, err := tarMetaData(meta, data)
	if err != nil {
		return nil, err
	}

	go func() { // ?
		err := s3.PutObject(backetStreamIn, uuid, t)
		if err != nil {
			// log.
		}
	}()

	return uuid, nil
}
