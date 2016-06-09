package core

import (
	"net/http"

	"internal/minio"
)

var Putd = &putdWorker{}

type putdWorker struct {
	meta []byte
	uuid string
}

func (w *putdWorker) ReadHeader(h http.Header) {
	w.meta = []byte(h.Get("Content-Meta"))
	w.uuid = h.Get("Content-UUID")
}

func (w *putdWorker) Work(data []byte) (interface{}, error) {
	t, err := tarMetaData(w.meta, data)
	if err != nil {
		return nil, err
	}

	go func() { // ?
		err := minio.PutObject(backetStreamIn, w.uuid, t)
		if err != nil {
			// log.
		}
	}()

	return w.uuid, nil
}
