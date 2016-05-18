package core

import (
	"net/http"

	"internal/redis"

	"golang.org/x/net/context"
)

// Ping calls Redis PING
func Ping(_ context.Context, _ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	c := redis.Get()
	defer redis.Put(c)

	return c.Do("PING")
}
