package main

import (
	"io/ioutil"
	"log"
	"os"

	_ "expvar"
	_ "net/http/pprof"

	"internal/api"
	"internal/minio"
	"internal/nats"
	"internal/pref"
	"internal/redis"
	"internal/server"
)

func main() {
	pref.Init()
	initLogger(pref.Verbose)

	err := initAndRun(
		pref.NATS,
		pref.MINIO,
		pref.REDIS,
		pref.SERVER,
	)
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

// Fail fast and explicit dependencies
func initAndRun(addrNATS, addrMINIO, addrREDIS, addrSERVER string) error {
	n, err := nats.Init(addrNATS)
	if err != nil {
		return err
	}

	m, err := minio.NewClient(addrMINIO)
	if err != nil {
		return err
	}

	r, err := redis.Init(addrREDIS)
	if err != nil {
		return err
	}

	h, err := api.Init(n, m, r)
	if err != nil {
		return err
	}

	return server.Run(addrSERVER, h)
}
