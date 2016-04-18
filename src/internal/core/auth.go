package core

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

// Redis scheme:
// SET key => val [val ...]

const keyAuth = "auth:list"

type (
	decodeAuth []byte

	// Redis scheme:
	// SET => key="auth"
	// SADD key v [v...]
	// SREM key v [v...]
	// SISMEMBER key v
)

func (d decodeAuth) src() ([]string, error) {
	var out []string
	err := json.Unmarshal(d, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (d decodeAuth) vls(withKey bool) ([]interface{}, error) {
	src, err := d.src()
	if err != nil {
		return nil, err
	}

	out := make([]interface{}, 0, len(src)+1)
	if withKey {
		out = append(out, keyAuth)
	}
	for i := range src {
		out = append(out, src[i])
	}

	return out, nil
}

func (d decodeAuth) get(c redis.Conn) ([]interface{}, error) {
	vls, err := d.vls(false)
	if err != nil {
		return nil, err
	}

	for i := range vls {
		err = c.Send("SISMEMBER", keyAuth, vls[i])
		if err != nil {
			return nil, err
		}
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	var rcv int64
	for i := range vls {
		rcv, err = redis.Int64(c.Receive())
		if err != nil {
			return nil, err
		}
		if rcv == 0 {
			vls[i] = nil
			continue
		}
		vls[i] = toString(vls[i])
	}

	return vls, nil
}

func (d decodeAuth) set(c redis.Conn) (interface{}, error) {
	vls, err := d.vls(true)
	if err != nil {
		return nil, err
	}

	return c.Do("SADD", vls...)
}

func (d decodeAuth) del(c redis.Conn) (interface{}, error) {
	vls, err := d.vls(true)
	if err != nil {
		return nil, err
	}

	return c.Do("SREM", vls...)
}
