package api

import (
	"log"
	"net/http"

	"internal/core"
	"internal/database/minio"
)

func putd(data []byte, r, _ http.Header) (interface{}, error) {
	meta := []byte(r.Get("Content-Meta"))

	m, err := core.UnmarshalMeta(meta)
	if err != nil {
		return nil, err
	}

	err = core.CheckHTag(m.HTag)
	if err != nil {
		return nil, err
	}

	go func() { // ?
		t, err := core.PackMetaData(meta, data)
		if err != nil {
			log.Println("putd: err: gztr:", err)
			return
		}

		f := core.MakeFileName(m.Auth.ID, m.UUID, core.ConvHTag(m.HTag))
		err = minio.Put(core.BucketStreamIn, f, t)
		if err != nil {
			log.Println("putd: err: save:", err)
		}
	}()

	return m.UUID, nil
}
