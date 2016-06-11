package core

import (
	"log"
	"net/http"

	"internal/gzpool"
	"internal/minio"
)

var Putd = newPutdWorker()

type putdWorker struct {
	meta []byte
	uuid string
}

func newPutdWorker() Worker {
	return &putdWorker{}
}

func (w *putdWorker) NewWorker() Worker {
	return newPutdWorker()
}

func (w *putdWorker) ReadHeader(h http.Header) {
	w.meta = []byte(h.Get("Content-Meta"))
	w.uuid = h.Get("Content-UUID")
}

func (w *putdWorker) Work(data []byte) (interface{}, error) {
	go func() { // ?
		data, err := gzpool.MustGzip(data)
		if err != nil {
			//return nil, err
		}

		t, err := tarMetaData(w.meta, data)
		if err != nil {
			//return nil, err
		}

		err = minio.PutObject(backetStreamIn, w.uuid, t)
		if err != nil {
			log.Println("putdWorker go func", err)
		}
	}()

	return w.uuid, nil
}
