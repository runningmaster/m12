package core

import (
	"encoding/base64"
	"log"
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
	h.Set("Content-Type", "gzip") // for writeResp
	h.Set("Content-Meta", base64.StdEncoding.EncodeToString(w.meta))
}

func (w *popdWorker) Work(data []byte) (interface{}, error) {
	p, err := unmarshaJSONpair(data)
	if err != nil {
		return nil, err
	}

	o, err := cMINIO.GetObject(p.Backet, p.Object)
	if err != nil {
		return nil, err
	}

	meta, data, err := untarMetaData(o)
	if err != nil {
		return nil, err
	}
	w.meta = meta

	go func(p pair) {
		err := cMINIO.RemoveObject(p.Backet, p.Object)
		if err != nil {
			log.Println("popMetaData", err)
		}
	}(p)

	return data, nil
}
