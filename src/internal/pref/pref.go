package pref

import (
	"expvar"
	"flag"
	"strconv"
)

// NOTE: about priority: default <- key/value store <- config <- env <- flag <-explicit set

var (
	// Verbose is flag for output
	Verbose = flag.Bool("verbose", true, "Verbose output")

	// Debug mode
	Debug = flag.Bool("debug", false, "Debug mode")

	// Masterkey is default secret key for sysdba
	Masterkey = flag.String("masterkey", "masterkey", "secret key for sysdba")
)

func init() {
	//experimentWithExpVarFIXME()
}

func experimentWithExpVarFIXME() {
	expvar.NewInt("debug").Set(1)
	d, _ := strconv.Atoi(expvar.Get("debug").String())
	*Debug = d == 1
}
