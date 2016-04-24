package main

import (
	_ "expvar"
	_ "net/http/pprof"

	"internal/flag"
	"internal/log"
	"internal/server"
)

func main() {
	flag.Parse()

	err := server.Run()
	if err != nil {
		log.Printf("main: %s", err)
	}
}
