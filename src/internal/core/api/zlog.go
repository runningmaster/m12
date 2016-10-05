package api

import (
	"internal/core/redis"
)

func getZlog(data []byte) (interface{}, error) {
	return redis.GetZlog(data)
}

func getMeta(data []byte) (interface{}, error) {
	return redis.GetMeta(data)
}
