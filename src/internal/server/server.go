package server

import (
	"time"

	"internal/api"
	"internal/flag"

	"github.com/tylerb/graceful"
)

func Run() error {
	err := api.Init(regFunc)
	if err != nil {
		return err
	}

	s, err := makeServer(flag.Addr)
	if err != nil {
		return err
	}

	return graceful.ListenAndServe(s, 5*time.Second)
}
