package api

import (
	"log"
	"net/http"

	"internal/core/structs"
	"internal/database/minio"
)

func putd(data []byte, r, _ http.Header) (interface{}, error) {
	meta := []byte(r.Get("Content-Meta"))

	m, err := structs.UnmarshalMeta(meta)
	if err != nil {
		return nil, err
	}

	err = structs.CheckHTag(m.HTag)
	if err != nil {
		return nil, err
	}

	go func() { // ?
		t, err := gztarMetaData(meta, data)
		if err != nil {
			log.Println("putd: err: gztr:", err)
			return
		}

		f := structs.MakeFileName(m.Auth.ID, m.UUID, structs.FindHTag(m.HTag))
		err = minio.Put(structs.BucketStreamIn, f, t)
		if err != nil {
			log.Println("putd: err: save:", err)
		}
	}()

	return m.UUID, nil
}
