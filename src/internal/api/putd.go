package api

import (
	"fmt"
	"log"
	"net/http"
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

		f := makeFileName(m.Auth.ID, m.UUID, convHTag[m.HTag])
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

func makeFileName(auth, uuid, htag string) string {
	return fmt.Sprintf("%s_%s_%s.tar", trimPart(auth), trimPart(uuid), htag)
}
