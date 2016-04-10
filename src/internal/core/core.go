package core

import (
	"net/http"

	"golang.org/x/net/context"
)

// Handler is func for processing data from api.
type Handler func(context.Context, *http.Request) (interface{}, error)
