package main

import (
	"io/ioutil"
	"log"
	"os"

	_ "expvar"
	_ "net/http/pprof"

	"internal/pref"
	"internal/server"
)

func main() {
	pref.Init()
	initLogger(pref.Verbose)

	err := server.Run(pref.Host)
	if err != nil {
		pref.Usage()
		log.Fatalf("main: %v", err)
	}
}

func initLogger(v bool) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	if v {
		log.SetOutput(os.Stderr)
	}
}
