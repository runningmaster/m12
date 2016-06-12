package pref

const envFormat = "M12_%s"

type pref struct {
	name  string
	usage string
	value interface{}
}

var (
	// Host is host server address.
	Host = "http://127.0.0.1:8080"

	// NATS is NATS server address.
	NATS = "nats://127.0.0.1:4222"

	// Minio is Minio server address.
	Minio = "http://akey:skey@127.0.0.1:9000"

	// Redis is Redis server address.
	Redis = "redis://127.0.0.1:6379"

	// MasterKey is default secret key for sysdba.
	MasterKey = "masterkey"

	// Debug is flag for debug mode.
	Debug = false

	// Verbose is flag for verbose output.
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
