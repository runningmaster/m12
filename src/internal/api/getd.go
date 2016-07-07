package api

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"sync"
)

var Getd = &getd{}

type getd struct {
	mu   sync.Mutex
	meta []byte
}

func (w *getd) New() *getd {
	w.mu.Lock()
	defer w.mu.Unlock()
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

	buf := new(bytes.Buffer)

	enc := base64.NewEncoder(base64.StdEncoding, buf)
	enc.Write(w.meta)
	enc.Close()
	h.Set("Content-Meta", buf.String())
	//h.Set("Content-Meta", base64.StdEncoding.EncodeToString(w.meta))
}
