package api

import (
	"encoding/base64"
	"net/http"
)

func getd(data []byte, _, w http.Header) (interface{}, error) {
	b, o, err := unmarshaPairExt(data)
	if err != nil {
		return nil, err
	}

	r, err := minio.Get(b, o)
	if err != nil {
		return nil, err
	}
	defer minio.Free(r)

	meta, data, err := ungztarMetaData(r, false, true)
	if err != nil {
		return nil, err
	}

	w.Set("Content-Encoding", "gzip")
	w.Set("Content-Type", "gzip") // for writeResp
	w.Set("Content-Meta", base64.StdEncoding.EncodeToString(meta))
	return data, nil
}
