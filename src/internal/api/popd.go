package api

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
)

func popd(data []byte, _, w http.Header) (interface{}, error) {
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

	w.Set("Content-Encoding", "gzip")
	w.Set("Content-Type", "gzip") // for writeResp
	w.Set("Content-Meta", base64.StdEncoding.EncodeToString(meta))
	return data, nil
}
