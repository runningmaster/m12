package main

import (
	_ "expvar"
	_ "net/http/pprof"

	_ "internal/api"
	"internal/flag"
	"internal/log"
	"internal/server"
	"internal/signal"
)

func main() {
	flag.Parse()

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
}
