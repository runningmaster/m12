package main

import (
	"io/ioutil"
	"log"
	"os"

	_ "expvar"
	_ "net/http/pprof"

	"internal/flag"
	"internal/server"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

func main() {
	flag.Parse()
	if flag.Verbose {
		log.SetOutput(os.Stderr)
	}

	err := server.Run()
	if err != nil {
		log.Printf("main: %s", err)
	}
}
