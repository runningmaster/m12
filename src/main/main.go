package main

import (
	"io/ioutil"
	"log"
	"os"
	//"runtime"

	_ "expvar"
	_ "net/http/pprof"

	//"internal/api"
	"internal/pref"
	"internal/server"
)

func main() {
	pref.Init()
	initLogger(pref.Verbose)

	err := server.Run(pref.Host)
	//err := api.TestStreamOut()
	if err != nil {
		log.Fatalf("main: %v", err)
	}
	//runtime.Goexit()
}

func initLogger(v bool) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	if v {
		log.SetOutput(os.Stderr)
	}
}
