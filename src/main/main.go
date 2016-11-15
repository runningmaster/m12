package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	_ "expvar"
	_ "net/http/pprof"

	"internal/core"
	"internal/core/api"
	"internal/core/pref"
	"internal/database/minio"
	"internal/database/redis"
	"internal/net/http/server"
	"internal/net/nats"
	"internal/version"
)

func main() {
	pref.Init()
	initLogger(systemdBased(), pref.Verbose)

	err := initAndRun(
		pref.NATS,
		pref.MINIO,
		pref.REDIS,
		pref.SERVER,
	)
	if err != nil {
		log.Println(version.AppName(), err)
		pref.Usage()
		os.Exit(1)
	}
}

func systemdBased() bool {
	return exec.Command("pidof", "systemd").Run() == nil
}

func initLogger(f, v bool) {
	log.SetOutput(ioutil.Discard)
	if f {
		log.SetFlags(0)
	}
	if v {
		log.SetOutput(os.Stderr)
	}
}

// Fail fast and explicit dependencies
func initAndRun(addrNATS, addrMINIO, addrREDIS, addrSERVER string) error {
	err := nats.Init(addrNATS)
	if err != nil {
		return err // we can receiving data any way or we can not (?)
	}

	err = minio.Init(addrMINIO)
	if err != nil {
		return err
	}

	err = redis.Init(addrREDIS)
	if err != nil {
		return err
	}

	err = core.Init()
	if err != nil {
		return err
	}

	return server.Run(addrSERVER, api.MakeRouter())
}
