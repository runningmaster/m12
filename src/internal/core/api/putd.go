package api

import (
	"log"
	"net/http"

	"internal/core/structs"
	"internal/database/minio"
)

func putd(data []byte, r, _ http.Header) (interface{}, error) {
	meta := []byte(r.Get("Content-Meta"))

	m, err := unmarshalMeta(meta)
	if err != nil {
		return nil, err
	}

	err = checkHTag(m.HTag)
	if err != nil {
		return nil, err
	}

	go func() { // ?
		t, err := gztarMetaData(meta, data)
		if err != nil {
			log.Println("putd: err: gztr:", err)
			return
		}

		f := structs.MakeFileName(m.Auth.ID, m.UUID, meta.FindHTag(m.HTag))
		err = minio.Put(bucketStreamIn, f, t)
		if err != nil {
			log.Println("putd: err: save:", err)
		}
	}()

	return m.UUID, nil
}
