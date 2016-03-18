package main

import (
	"internal/log"
	"internal/server"
	"internal/signal"
	"internal/version"
)

import (
	_ "expvar"
	_ "net/http/pprof"

	_ "internal/api"
)

func main() {
	log.Printf("main: start version %s", version.Stamp)

	errCh := make(chan error)
	go func(ch chan<- error) {
		ch <- server.Start()
		close(ch)
	}(errCh)

	select {
	case err := <-errCh:
		log.Printf("main: error occurred: %s", err)
	case sig := <-signal.WaitToExit():
		log.Printf("main: signal received: %s", sig)
		err := server.Stop()
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("main: the end")
}
