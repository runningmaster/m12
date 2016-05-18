package core

import "github.com/garyburd/redigo/redis"

var redisServer connGetter

type connGetter interface {
	Get() redis.Conn
}

type redisGetSetDelOper interface {
	get(redis.Conn) ([]interface{}, error)
	set(redis.Conn) (interface{}, error)
	del(redis.Conn) (interface{}, error)
}
