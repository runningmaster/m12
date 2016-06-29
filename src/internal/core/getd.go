package core

import (
	"encoding/base64"
	"io"
	"net/http"
)

var Getd = &getd{}

type getd struct {
	meta []byte
}

func (w *getd) New() interface{} {
	return &getd{}
}

func (w *getd) Work(data []byte) (interface{}, error) {
	bucket, object, err := unmarshaPairExt(data)
	if err != nil {
		return nil, err
	}

	o, err := cMINIO.GetObject(bucket, object)
	if err != nil {
		return nil, err
	}

	defer func(c io.Closer) {
		if c != nil {
			_ = c.Close()
		}
	}(o)

	meta, data, err := ungztarMetaData(o, false, true)
	if err != nil {
		return nil, err
	}
	w.meta = meta

	return data, nil
}

func (w *getd) WriteHeader(h http.Header) {
	h.Set("Content-Encoding", "gzip")
	h.Set("Content-Type", "gzip") // for writeResp
	h.Set("Content-Meta", base64.StdEncoding.EncodeToString(w.meta))
}
