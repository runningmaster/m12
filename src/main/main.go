package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	_ "expvar"
	_ "net/http/pprof"

	"internal/api"
	"internal/nats"
	"internal/pref"
	"internal/redis"
	"internal/s3"
	"internal/server"
)

var (
	flagS3     = *flag.String("s3", "127.0.0.1:9000", "S3 object storage address")
	flagS3ak   = *flag.String("s3ak", "", "S3 access key")
	flagS3sk   = *flag.String("s3sk", "", "S3 secret key")
	flagNATS   = *flag.String("nats", "nats://user:pass@ip:4222", "network address for NATS server 'scheme://[user:pass]@host[:port]'")
	flagRedis  = *flag.String("redis", "redis://127.0.0.1:6379", "network address for Redis server 'scheme://[user:pass]@host[:port]'")
	flagServer = *flag.String("server", "127.0.0.1:8080", "server address '[host]:port'")
)

func main() {
	flag.Parse()
	l := makeLogger(pref.Verbose)

	err := initDepend(l)
	if err != nil {
		log.Fatalf("main: %v", err)
	}
}

func initDepend(l *log.Logger) error {
	err := s3.Run(flagS3, flagS3ak, flagS3sk, l)
	if err != nil {
		return err
	}

	err = nats.Run(flagNATS, l)
	if err != nil {
		return err
	}

	err = redis.Run(flagRedis, l)
	if err != nil {
		return err
	}

	err = api.Reg()
	if err != nil {
		return err
	}

	return server.Run(flagServer, nil)
}

func makeLogger(v bool) *log.Logger {
	out := ioutil.Discard
	if v {
		out = os.Stderr
	}
	return log.New(out, "", 0)
}
