package pref

import (
	"expvar"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const envFormat = "M12_%s"

type pref struct {
	name  string
	usage string
	value interface{}
}

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

	prefs = []pref{
		pref{
			"host",
			"Host server address",
			&Host,
		},
		pref{
			"nats",
			"NATS server address",
			&NATS,
		},
		pref{
			"minio",
			"Minio S3 object storage address",
			&Minio,
		},
		pref{
			"minio-akey",
			"Minio S3 access key",
			&MinioAKey,
		},
		pref{
			"minio-skey",
			"Minio S3 secret key",
			&MinioSKey,
		},
		pref{
			"redis",
			"Redis server address",
			&Redis,
		},
		pref{
			"masterkey",
			"Secret key for sysdba",
			&MasterKey,
		},
		pref{
			"debug",
			"Debug mode",
			&Debug,
		},
		pref{
			"verbose",
			"Verbose output",
			&Verbose,
		},
	}
)

// Init is public init func, must be called from main()
// NOTE about priority levels:
// 5.explicit set > 4.flag > 3.environment > 2.config > 1.key/value store > 0.default
func Init() {
	init0Defaults()
	init1KeyValStore()
	init2ConfigFile()
	init3FromEvars()
	init4FromFlags()
	init5ExpvarFIXME("debug")
}

func init0Defaults() {
	// see above
}

func init1KeyValStore() {
	// not implemented
}

func init2ConfigFile() {
	// not implemented
}

func init3FromEvars() {
	for i := range prefs {
		prefs[i].setFromEvar()
	}
}

func init4FromFlags() {
	for i := range prefs {
		prefs[i].setFromFlag()
	}
	flag.Parse()
}

func init5ExpvarFIXME(key string) {
	expvar.NewInt(key).Set(1)
	d, _ := strconv.Atoi(expvar.Get(key).String())
	Debug = d == 1
}

func (p pref) setFromFlag() {
	switch x := p.value.(type) {
	case *string:
		flag.StringVar(x, p.name, *x, p.usage)
	case *bool:
		flag.BoolVar(x, p.name, *x, p.usage)
	default:
		panic("pref: unreachable, add new flag case")
	}
}

func (p pref) envVar() string {
	return fmt.Sprintf(envFormat, strings.ToUpper(strings.Replace(p.name, "-", "_", -1)))
}

func (p pref) setFromEvar() {
	v := os.Getenv(p.envVar())
	if v == "" {
		return
	}

	switch x := p.value.(type) {
	case *string:
		*x = v
	case *bool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		*x = b
	default:
		panic("pref: unreachable, add new evar case")
	}
}
