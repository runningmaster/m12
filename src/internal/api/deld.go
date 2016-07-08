package api

import (
	"net/http"
)

func deld(data []byte, _, _ http.Header) (interface{}, error) {
	bucket, object, err := unmarshaPairExt(data)
	if err != nil {
		return nil, err
	}

	err = cMINIO.RemoveObject(bucket, object)
	if err != nil {
		return nil, err
	}

	return "OK", nil
}
