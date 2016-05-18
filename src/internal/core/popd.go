package core

import (
	"net/http"

	"internal/s3"

	"golang.org/x/net/context"
)

// Popd returns data from stream queues
func Popd(_ context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	b, err := readClose(r.Body)
	if err != nil {
		return nil, err
	}

	o, err := s3.PopObjectByPathJSON(b)
	if err != nil {
		return nil, err
	}

	meta, data, err := untarMetaData(o)
	if err != nil {
		return nil, err
	}

	m, err := makeMetaFromJSON(meta)
	if err != nil {
		return nil, err
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Meta", m.packToBase64String()) // FIXME

	return data, nil
}
