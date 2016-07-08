package pipe

import (
	"net/http"
)

type handler func(h http.Handler) http.Handler

func Use(pipes ...handler) http.Handler {
	var h http.Handler
	for i := len(pipes) - 1; i >= 0; i-- {
		h = pipes[i](h) // note: nill will generate panic
	}
	return h
}
