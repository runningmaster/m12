package core

import (
	"net/http"

	"internal/database/redispool"

	"golang.org/x/net/context"
)

// Ping calls Redis PING
func Ping(_ context.Context, _ *http.Request) (interface{}, error) {
	c := redispool.Get()
	defer redispool.Put(c)

	return c.Do("PING")
}
