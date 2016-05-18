package main

import (
	"io/ioutil"
	"log"
	"os"

	_ "expvar"
	_ "net/http/pprof"

	"internal/conf"
	"internal/server"
)

func init() {
	initConfig()
	initLogger()
}

func main() {
	err := server.Run(conf.HostAddr)
	if err != nil {
		log.Printf("main: %s", err)
	}
}

func initConfig() {
	conf.Parse()
}

func initLogger() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	if conf.Verbose {
		log.SetOutput(os.Stderr)
	}
}
