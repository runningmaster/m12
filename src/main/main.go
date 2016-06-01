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
	flagMinio     = *flag.String("minio", "minio://127.0.0.1:9000", "Minio (S3) object storage address")
	flagMinioAKey = *flag.String("minio-akey", "", "Minio (S3) access key")
	flagMinioSKey = *flag.String("minio-skey", "", "Minio (S3) secret key")
	flagNATS      = *flag.String("nats", "nats://user:pass@ip:4222", "network address for NATS server 'scheme://[user:pass]@host[:port]'")
	flagRedis     = *flag.String("redis", "redis://127.0.0.1:6379", "network address for Redis server 'scheme://[user:pass]@host[:port]'")
	flagServer    = *flag.String("server", "http://127.0.0.1:8080", "server address '[host]:port'")
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
	err := s3.Run(flagMinio, flagMinioAKey, flagMinioSKey, l)
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
