package main

import (
	"flag"
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

var (
	flagNATS      = flag.String("nats", "nats://user:pass@host:4222", "NATS server address")
	flagMinio     = flag.String("minio", "http://127.0.0.1:9000", "Minio S3 object storage address")
	flagMinioAKey = flag.String("minio-akey", "", "Minio S3 access key")
	flagMinioSKey = flag.String("minio-skey", "", "Minio S3 secret key")
	flagRedis     = flag.String("redis", "redis://127.0.0.1:6379", "Redis server address")
	flagHost      = flag.String("host", "http://127.0.0.1:8080", "Host server address")
)

func main() {
	flag.Parse()
	l := makeLogger(*pref.Verbose)

	err := initDepend(l)
	if err != nil {
		l.Fatalf("main: %v", err)
	}
}

func initDepend(l *log.Logger) error {
	err := nats.Run(*flagNATS, l)
	if err != nil {
		return err
	}

	err = minio.Run(*flagMinio, *flagMinioAKey, *flagMinioSKey, l)
	if err != nil {
		return err
	}

	err = redis.Run(*flagRedis, l)
	if err != nil {
		return err
	}

	err = api.Reg()
	if err != nil {
		return err
	}

	return server.Run(*flagHost, nil)
}

func makeLogger(v bool) *log.Logger {
	out := ioutil.Discard
	if v {
		out = os.Stderr
	}
	return log.New(out, "", 0)
}
