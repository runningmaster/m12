package core

import (
	"net/http"

	"golang.org/x/net/context"
)

// ErrFailFast informs us that work is impossible
var ErrFailFast error

// Handler is func for processing data from api.
type Handler func(context.Context, *http.Request) (interface{}, error)
