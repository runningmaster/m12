package core

import (
	"net/http"

	"internal/redis"

	"golang.org/x/net/context"
)

// Ping calls Redis PING
func Ping(_ context.Context, _ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	return redis.Ping()
}
