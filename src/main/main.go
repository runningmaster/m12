package main

import (
	"io/ioutil"
	"log"
	"os"

	_ "expvar"
	_ "net/http/pprof"

	"internal/conf"
	//	"internal/nats"
	//	"internal/redis"
	//	"internal/s3"
	"internal/server"
)

func init() {
	initConfig()
	initLogger()
}

func main() {
	var err error

	//	err = s3.Run(conf.S3Address, conf.S3AccessKey, conf.S3SecretKey, nil)
	//	if err != nil {
	//		goto fail
	//	}

	//	err = nats.Run(conf.NATSAddress, nil)
	//	if err != nil {
	//		goto fail
	//	}

	//	err = redis.Run(conf.RedisAddress, nil)
	//	if err != nil {
	//		goto fail
	//	}

	err = server.Run(conf.ServerAddress)
	if err != nil {
		goto fail
	}

fail:
	if err != nil { // workaround for Ctrl+C
		log.Fatalf("main: %v", err)
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
