package flag

import (
	"expvar"
	"flag"
	"strconv"
)

// CtxKey is type for context keys
type CtxKey int

const (
	// CtxUUID is key for context value
	CtxUUID CtxKey = iota
	// CtxAuth is key for context value
	CtxAuth
	// CtxFail is key for context value
	CtxFail
	// CtxSize is key for context value
	CtxSize
	// CtxCode is key for context value
	CtxCode
	// CtxTime is key for context value
	CtxTime
)

// NOTE: about priority: default <- key/value store <- config <- env <- flag <-explicit set

var (
	// Addr is HTTP service address
	Addr = *flag.String("addr", "127.0.0.1:8080", "HTTP service address '[host]:port'")

	// Redis helps to assign a TCP network address for Redis server
	Redis = *flag.String("redis", "redis://127.0.0.1:6379", "A TCP network address for Redis server 'scheme://[user:pass]@host[:port]'")

	// Verbose is flag for output
	Verbose = *flag.Bool("verbose", true, "Verbose output")

	// Debug mode
	Debug = *flag.Bool("debug", false, "Debug mode")

	// Masterkey is default secret key for sysdba
	Masterkey = *flag.String("masterkey", "masterkey", "secret key for sysdba")

	// S3Address
	S3Address = *flag.String("s3addr", "127.0.0.1:9000", "S3 object storage address")

	// S3AccessKey
	S3AccessKey = *flag.String("s3akey", "ACCESSKEYID", "S3 access key")

	// S3SecretKey
	S3SecretKey = *flag.String("s3skey", "SECRETACCESSKEY", "S3 secret key")
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
