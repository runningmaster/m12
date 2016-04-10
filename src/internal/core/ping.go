package core

import (
	"net/http"

	"golang.org/x/net/context"
)

// Ping calls Redis PING
func Ping(_ context.Context, _ *http.Request) (interface{}, error) {
	c := redisPool.Get()
	defer redisPool.Put(c)

	return c.Do("PING")
}
