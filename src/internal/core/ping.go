package core

import (
	"net/http"

	"golang.org/x/net/context"
)

// Ping calls Redis PING
func Ping(_ context.Context, _ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	c := redisGet()
	defer func() { _ = redisPut(c) }()

	return c.Do("PING")
}
