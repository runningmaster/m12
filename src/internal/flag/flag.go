package flag

import (
	"expvar"
	"flag"
	"strconv"
)

// NOTE about priority: default <- key/value store <- config <- env <- flag <-explicit set

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
)

func init() {
	flag.Parse()
	experimentWithExpVarFIXME()
}

func experimentWithExpVarFIXME() {
	expvar.NewInt("debug").Set(1)
	d, _ := strconv.Atoi(expvar.Get("debug").String())
	Debug = d == 1
}
