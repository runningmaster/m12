package core

import (
	"fmt"
	"log"
	"net/http"

	"internal/gzpool"
)

var Putd = newPutdWorker()

type putdWorker struct {
	meta []byte
}

func newPutdWorker() Worker {
	return &putdWorker{}
}

func (w *putdWorker) NewWorker() Worker {
	return newPutdWorker()
}

func (w *putdWorker) ReadHeader(h http.Header) {
	w.meta = []byte(h.Get("Content-Meta"))
}

func (w *putdWorker) Work(data []byte) (interface{}, error) {
	m, err := unmarshalJSONmeta(w.meta)
	if err != nil {
		return nil, err
	}

	err = checkHTag(m.HTag)
	if err != nil {
		return nil, err
	}

	go func() { // ?
		d, err := gzpool.MustGzip(data)
		if err != nil {
			log.Println("putd: gzip:", err)
		}

		t, err := tarMetaData(w.meta, d)
		if err != nil {
			log.Println("putd: tar:", err)
		}

		f := makeFileName(m.UUID, m.Auth, m.HTag)
		_, err = cMINIO.PutObject(backetStreamIn, f, t, "")
		if err != nil {
			log.Println("putd: minio: err:", err)
		}
	}()

	return m.UUID, nil
}

const magicLen = 8

func makeFileName(uuid, auth, htag string) string {
	return fmt.Sprintf("%s_%s_%s.tar", uuid, auth[:magicLen], htag)
}
