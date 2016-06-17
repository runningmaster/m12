package core

import (
	"fmt"
	"log"
	"net/http"
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
	m, err := unmarshalMeta(w.meta)
	if err != nil {
		return nil, err
	}

	err = checkHTag(m.HTag)
	if err != nil {
		return nil, err
	}

	go func() { // ?
		t, err := gztarMetaData(w.meta, data)
		if err != nil {
			log.Println("putd: err: gztr:", err)
			return
		}

		f := makeFileName(m.UUID, m.Auth, m.HTag)
		_, err = cMINIO.PutObject(bucketStreamIn, f, t, "")
		if err != nil {
			log.Println("putd: err: save:", err)
		}
	}()

	return m.UUID, nil
}

const magicLen = 8

func trimPart(s string) string {
	if len(s) > magicLen {
		return s[:magicLen]
	}
	return s
}

func makeFileName(uuid, auth, htag string) string {
	return fmt.Sprintf("%s_%s_%s.tar", trimPart(uuid), trimPart(auth), htag)
}
