package pref

import (
	"expvar"
	"flag"
	"strconv"
)

// NOTE about priority levels:
// 5.explicit set > 4.flag > 3.environment > 2.config > 1.key/value store > 0.default

var (
	// Host is host server address
	Host = "http://127.0.0.1:8080"

	// NATS is NATS server address
	NATS = "nats://user:pass@host:4222"

	// Minio is Minio server address
	Minio = "http://127.0.0.1:9000"

	// MinioAKey is Minio access key
	MinioAKey = ""

	// MinioSKey is Minio secret key
	MinioSKey = ""

	// Redis is Redis server address
	Redis = "redis://127.0.0.1:6379"

	// MasterKey is default secret key for sysdba
	MasterKey = "masterkey"

	// Debug is flag for debug mode
	Debug = false

	// Verbose is flag for verbose output
	Verbose = true
)

func init0() {
	// see above
}

func init1() {
	// not implemented
}

func init2() {
	// not implemented
}

func init3() {
	// not implemented
}

func init4() {
	flagString("host", "Host server address", &Host)
	flagString("nats", "NATS server address", &NATS)
	flagString("minio", "Minio S3 object storage address", &Minio)
	flagString("minio-akey", "Minio S3 access key", &MinioAKey)
	flagString("minio-skey", "Minio S3 secret key", &MinioSKey)
	flagString("redis", "Redis server address", &Redis)
	flagString("masterkey", "Secret key for sysdba", &MasterKey)
	flagBool("debug", "Debug mode", &Debug)
	flagBool("verbose", "Verbose output", &Verbose)
	flag.Parse()
}

func init5() {
	// not implemented
	// expvar (?)
}

func Init() {
	init0()
	init1()
	init2()
	init3()
	init4()
	init5()
}

func flagString(name, usage string, val *string) {
	flag.StringVar(val, name, *val, usage)
}

func flagBool(name, usage string, val *bool) {
	flag.BoolVar(val, name, *val, usage)
}

func experimentWithExpVarFIXME() {
	expvar.NewInt("debug").Set(1)
	d, _ := strconv.Atoi(expvar.Get("debug").String())
	Debug = d == 1
}
