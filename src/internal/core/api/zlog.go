package api

import (
	"internal/core"
)

func getZlog(data []byte) (interface{}, error) {
	return core.GetZlog(data)
}

func getMeta(data []byte) (interface{}, error) {
	return core.GetMeta(data)
}
