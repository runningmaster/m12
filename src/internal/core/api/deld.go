package api

import (
	"internal/core"
	"internal/database/minio"
)

func deld(data []byte) (interface{}, error) {
	b, o, err := core.DecodePath(data)
	if err != nil {
		return nil, err
	}

	err = minio.Del(b, o)
	if err != nil {
		return nil, err
	}

	return "OK", nil
}
