package server

import (
	"fmt"
	"sync"

	"internal/flag"

	"github.com/braintree/manners"
)

var (
	// FailFast informs us that work is impossible
	FailFast error

	once sync.Once
	gsrv *manners.GracefulServer
)

// initOnce sets once *http.Server for Start/Stop
func initOnce() error {
	var errOnce error
	once.Do(func() {
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

	if errOnce != nil {
		return errOnce
	}

	return nil
}

// Start starts HTTP server
func Start() error {
	if FailFast != nil {
		return fmt.Errorf("server: fail fast: %v", FailFast)
	}

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
