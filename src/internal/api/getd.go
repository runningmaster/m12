package api

import (
	"encoding/base64"
	"net/http"
)

func getd(data []byte, _, w http.Header) (interface{}, error) {
	b, o, err := cMINIO.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	f, err := cMINIO.Get(b, o)
	if err != nil {
		return nil, err
	}
	defer cMINIO.Free(f)

	m, d, err := ungztarMetaData(f, false, true)
	if err != nil {
		return nil, err
	}

	w.Set("Content-Encoding", "gzip")
	w.Set("Content-Type", "gzip") // for writeResp
	w.Set("Content-Meta", base64.StdEncoding.EncodeToString(m))
	return d, nil
}
