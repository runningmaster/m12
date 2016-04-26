package core

import (
	"net/http"

	"golang.org/x/net/context"
)

// Popd returns data from stream queues
func Popd(_ context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	b, err := readClose(r.Body)
	if err != nil {
		return nil, err
	}

	p, err := makePairFromJSON(b)
	if err != nil {
		return nil, err
	}

	meta, data, err := mineMetaData(p.Backet, p.Object)
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

func mineMetaData(backet, object string) ([]byte, []byte, error) {
	o, err := popObject(backet, object)
	if err != nil {
		return nil, nil, err
	}

	return untarMetaData(o)
}
