package core

import (
	"net/http"

	"golang.org/x/net/context"
)

// Handler is func for processing data from api.
type Handler func(context.Context, *http.Request) (interface{}, error)

func Init() error {
	var err error
	if err = initRedis(); err != nil {
		return err
	}

	if err = initCliS3(); err != nil {
		return err
	}

	return nil
}
