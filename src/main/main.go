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

	err := nats.Run(pref.NATS)
	if err != nil {
		failFast(err)
	}

	err = minio.Run(pref.Minio, pref.MinioAKey, pref.MinioSKey)
	if err != nil {
		failFast(err)
	}

	err = redis.Run(pref.Redis)
	if err != nil {
		failFast(err)
	}

	err = api.Reg()
	if err != nil {
		failFast(err)
	}

	err = server.Run(pref.Host)
	if err != nil {
		failFast(err)
	}
}

func initLogger(v bool) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	if v {
		log.SetOutput(os.Stderr)
	}
}

func failFast(err error) {
	log.Fatalf("main: %v", err)
}
