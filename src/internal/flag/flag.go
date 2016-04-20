package flag

import (
	"expvar"
	"flag"
	"strconv"
)

// NOTE: about priority: default <- key/value store <- config <- env <- flag <-explicit set

var (
	// Addr is HTTP service address
	Addr = *flag.String("addr", "127.0.0.1:8080", "service address '[host]:port'")

	// Redis is a TCP network address for Redis server
	Redis = *flag.String("redis", "redis://127.0.0.1:6379", "network address for Redis server 'scheme://[user:pass]@host[:port]'")

	// Verbose is flag for output
	Verbose = *flag.Bool("verbose", true, "Verbose output")

	// Debug mode
	Debug = *flag.Bool("debug", false, "Debug mode")

	// Masterkey is default secret key for sysdba
	Masterkey = *flag.String("masterkey", "masterkey", "secret key for sysdba")

	// S3Address is S3 object storage address
	S3Address = *flag.String("s3addr", "127.0.0.1:9000", "S3 object storage address")

	// S3AccessKey is S3 access key
	S3AccessKey = *flag.String("s3akey", "WTTDDZJKU1AJGX1CNNDD", "S3 access key")

	// S3SecretKey is S3 secret key
	S3SecretKey = *flag.String("s3skey", "jqBA9Qlzhn84ArN9s9AnwRQr3mkjlLz713tlHxCH", "S3 secret key")

	// NATS is a TCP network address for NATS server
	NATS = *flag.String("nats", "nats://@:4222", "network address for NATS server 'scheme://[user:pass]@host[:port]'")
)

// Parse is wrapper for std flag.Parse()
func Parse() {
	flag.Parse()
	experimentWithExpVarFIXME()
}

func experimentWithExpVarFIXME() {
	expvar.NewInt("debug").Set(1)
	d, _ := strconv.Atoi(expvar.Get("debug").String())
	Debug = d == 1
}
