package core

import (
	"encoding/base64"
	"net/http"
)

var Popd = newPopdWorker()

type popdWorker struct {
	meta []byte
}

func newPopdWorker() Worker {
	return &popdWorker{}
}

func (w *popdWorker) NewWorker() Worker {
	return newPopdWorker()
}

func (w *popdWorker) WriteHeader(h http.Header) {
	h.Set("Content-Encoding", "gzip")
	h.Set("Content-Meta", base64.StdEncoding.EncodeToString(w.meta))
}

func (w *popdWorker) Work(data []byte) (interface{}, error) {
	meta, data, err := popMetaData(data)
	if err != nil {
		return nil, err
	}
	w.meta = meta

	return data, nil
}
