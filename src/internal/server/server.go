package server

import (
	"time"

	"internal/api"

	"github.com/tylerb/graceful"
)

func Run(addr string) error {
	err := api.Init(regFunc)
	if err != nil {
		return err
	}

	s, err := makeServer(addr)
	if err != nil {
		return err
	}

	return graceful.ListenAndServe(s, 5*time.Second)
}
