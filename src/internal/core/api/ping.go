package api

import (
	"internal/core"
)

func ping() (interface{}, error) {
	return core.Ping()
}

func info() (interface{}, error) {
	return core.Info()
}
