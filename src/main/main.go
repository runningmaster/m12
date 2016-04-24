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

func main() {
	flag.Parse()
	initLogger()
	execServer()
}

func initLogger() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	if flag.Verbose {
		log.SetOutput(os.Stderr)
	}
}

func execServer() {
	err := server.Run()
	if err != nil {
		log.Printf("main: %s", err)
	}
}
