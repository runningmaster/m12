package server

import (
	"sync"

	"internal/errors"
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
		return errors.Locus(errOnce)
	}

	return nil
}

// Start starts HTTP server
func Start() error {
	if FailFast != nil {
		return errors.Locusf("server: fail fast: %v", FailFast)
	}

	if err := initOnce(); err != nil {
		return errors.Locus(err)
	}

	return gsrv.ListenAndServe()
}

// Stop stops HTTP server
func Stop() error {
	if gsrv == nil {
		return errors.Locusf("server: not registered")
	}

	if ok := gsrv.Close(); !ok {
		return errors.Locusf("server: received false on close")
	}

	return nil
}
