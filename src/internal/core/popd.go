package core

import (
	"encoding/base64"
	"io"
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

func (w *popdWorker) Work(data []byte) (interface{}, error) {
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

	defer func() {
		err = cMINIO.RemoveObject(bucket, object)
		if err != nil {
			log.Println("minio:", object, err)
		}
	}()

	meta, data, err := ungztarMetaData(o, false, true)
	if err != nil {
		return nil, err
	}
	w.meta = meta

	return data, nil
}

func (w *popdWorker) WriteHeader(h http.Header) {
	h.Set("Content-Encoding", "gzip")
	h.Set("Content-Type", "gzip") // for writeResp
	h.Set("Content-Meta", base64.StdEncoding.EncodeToString(w.meta))
}
