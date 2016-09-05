package pref

import "flag"

const envFormat = "M12_%s"

type pref struct {
	name  string
	usage string
	value interface{}
}

var (
	// NATS is NATS server address.
	NATS = "nats://127.0.0.1:4222"

	// MINIO is Minio server address.
	MINIO = "http://akey:skey@127.0.0.1:9000"

	// REDIS is Redis server address.
	REDIS = "redis://127.0.0.1:6379"

	// SERVER is host server address.
	SERVER = "http://127.0.0.1:8080"

	// MasterKey is default secret key for sysdba.
	MasterKey = "masterkey"

	// Debug is flag for debug mode.
	Debug = false

	// Verbose is flag for verbose output.
	Verbose = true

	prefs = []pref{
		pref{
			"nats",
			"NATS server address",
			&NATS,
		},
		pref{
			"minio",
			"Minio S3 object storage address",
			&MINIO,
		},
		pref{
			"redis",
			"Redis server address",
			&REDIS,
		},
		pref{
			"host",
			"Host server address",
			&SERVER,
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
	init0(prefs...) // 0 default
	init1(prefs...) // 1 key/value store
	init2(prefs...) // 2 config
	init3(prefs...) // 3 environment
	init4(prefs...) // 4 flag
	init5("debug")  // 5 explicit set TEST ONLY, FIXME
}

// Usage wraps flag.Usage
func Usage() {
	flag.Usage()
}
