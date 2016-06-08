package core

import (
	"bytes"
	"net/http"

	"internal/minio"
)

var Putd = &putdWorker{}

type putdWorker struct {
	meta []byte
}

func (w *putdWorker) ReadHeader(h http.Header) {
	w.meta = []byte(h.Get("Content-Meta"))
}

func (w *putdWorker) Work(data []byte) (interface{}, error) {
	pos := bytes.IndexByte(w.meta, '.')
	uuid := string(w.meta[:pos])
	meta := w.meta[pos+1:]

	t, err := tarMetaData(meta, data)
	if err != nil {
		return nil, err
	}

	go func() { // ?
		err := minio.PutObject(backetStreamIn, uuid, t)
		if err != nil {
			// log.
		}
	}()

	return uuid, nil
}
