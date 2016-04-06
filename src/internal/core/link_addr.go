package core

import (
	"encoding/json"

	"internal/errors"

	"github.com/garyburd/redigo/redis"
)

// Redis scheme:
// HASH key => [l/v] [a/v] [s/v] [e/v]

type decodeLinkAddr []byte

func (d decodeLinkAddr) src() ([]string, error) {
	var out []string
	if err := json.Unmarshal(d, &out); err != nil {
		return nil, errors.Locus(err)
	}

	return out, nil
}

func (d decodeLinkAddr) lnk() ([]linkAddr, error) {
	var out []linkAddr
	if err := json.Unmarshal(d, &out); err != nil {
		return nil, errors.Locus(err)
	}

	return out, nil
}

func (d decodeLinkAddr) vls() ([]interface{}, error) {
	src, err := d.src()
	if err != nil {
		return nil, errors.Locus(err)
	}

	out := make([]interface{}, 0, len(src))
	for i := range src {
		out = append(out, src[i])
	}

	return out, nil
}

func (d decodeLinkAddr) get(c redis.Conn) ([]interface{}, error) {
	src, err := d.src()
	if err != nil {
		return nil, errors.Locus(err)
	}

	var l linkAddr
	for i := range src {
		if err = c.Send("HMGET", l.keyflds(src[i])...); err != nil {
			return nil, err
		}
	}

	if err = c.Flush(); err != nil {
		return nil, errors.Locus(err)
	}

	out := make([]interface{}, 0, len(src))
	var rcv []interface{}
	for i := range src {
		if rcv, err = redis.Values(c.Receive()); err != nil {
			return nil, errors.Locus(err)
		}
		out = append(out, l.makeFrom(src[i], rcv))
	}

	return out, nil
}

func (d decodeLinkAddr) set(c redis.Conn) (interface{}, error) {
	lnk, err := d.lnk()
	if err != nil {
		return nil, err
	}

	for i := range lnk {
		if err = c.Send("DEL", lnk[i].ID); err != nil {
			return nil, errors.Locus(err)
		}
		if err = c.Send("HMSET", lnk[i].keyvals()...); err != nil {
			return nil, errors.Locus(err)
		}
	}

	if err = c.Flush(); err != nil {
		return nil, errors.Locus(err)
	}

	for range lnk {
		if _, err = c.Receive(); err != nil {
			return nil, errors.Locus(err)
		}
	}

	return stringOK(), nil
}

func (d decodeLinkAddr) del(c redis.Conn) (interface{}, error) {
	vls, err := d.vls()
	if err != nil {
		return nil, errors.Locus(err)
	}

	return c.Do("DEL", vls...)
}
