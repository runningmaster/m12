package server

import (
	"fmt"
	"sync"

	"internal/api"
	"internal/flag"

	"github.com/braintree/manners"
)

var (
	once sync.Once
	gsrv *manners.GracefulServer
)

// initOnce sets once *http.Server for Start/Stop
func initOnce() error {
	var errOnce error
	once.Do(func() {
		if err := api.Init(regFunc); err != nil {
			errOnce = err
			return
		}

		r, err := makeRouter()
		if err != nil {
			errOnce = err
			return
		}

		s, err := withRouter(flag.Addr, r)
		if err != nil {
			errOnce = err
			return
		}

		gsrv = manners.NewWithServer(s)
	})

	return errOnce
}

// Start starts HTTP server
func Start() error {
	if err := initOnce(); err != nil {
		return err
	}

	return gsrv.ListenAndServe()
}

// Stop stops HTTP server
func Stop() error {
	if gsrv == nil {
		return fmt.Errorf("server: not registered")
	}

	if ok := gsrv.Close(); !ok {
		return fmt.Errorf("server: received false on close")
	}

	return nil
}
