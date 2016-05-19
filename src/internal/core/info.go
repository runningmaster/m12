package core

import (
	"net/http"

	"internal/redis"

	"golang.org/x/net/context"
)

// Info calls Redis INFO
func Info(_ context.Context, _ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	return redis.Info()
}
