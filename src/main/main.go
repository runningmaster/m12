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
	l := makeLogger(pref.Verbose)

	var err error
	err = nats.Run(pref.NATS, l)
	if err != nil {
		goto fail
	}

	err = minio.Run(pref.Minio, pref.MinioAKey, pref.MinioSKey, l)
	if err != nil {
		goto fail
	}

	err = redis.Run(pref.Redis, l)
	if err != nil {
		goto fail
	}

	err = api.Reg()
	if err != nil {
		goto fail
	}

	err = server.Run(pref.Host, nil)
	if err != nil {
		goto fail
	}

fail:
	if err != nil {
		l.Fatalf("main: %v", err)
	}
}

func makeLogger(v bool) *log.Logger {
	out := ioutil.Discard
	if v {
		out = os.Stderr
	}
	return log.New(out, "", 0)
}
