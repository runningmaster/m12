package core

import (
	"encoding/json"

	"internal/errors"

	"github.com/garyburd/redigo/redis"
)

// HASH key => [l/v] [d/v] [b/v] [c/v] [s/v]

type decodeLinkDrug []byte

func (d decodeLinkDrug) src() ([]string, error) {
	var out []string
	if err := json.Unmarshal(d, &out); err != nil {
		return nil, errors.Locus(err)
	}

	return out, nil
}

func (d decodeLinkDrug) lnk() ([]linkDrug, error) {
	var out []linkDrug
	if err := json.Unmarshal(d, &out); err != nil {
		return nil, errors.Locus(err)
	}

	return out, nil
}

func (d decodeLinkDrug) vls() ([]interface{}, error) {
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

func (d decodeLinkDrug) get(c redis.Conn) ([]interface{}, error) {
	src, err := d.src()
	if err != nil {
		return nil, errors.Locus(err)
	}

	var l linkDrug
	for i := range src {
		if err = c.Send("HMGET", l.keyflds(src[i])...); err != nil {
			return nil, errors.Locus(err)
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

func (d decodeLinkDrug) set(c redis.Conn) (interface{}, error) {
	lnk, err := d.lnk()
	if err != nil {
		return nil, errors.Locus(err)
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

	for _ = range lnk {
		if _, err = c.Receive(); err != nil {
			return nil, errors.Locus(err)
		}
	}

	return "OK", nil
}

func (d decodeLinkDrug) del(c redis.Conn) (interface{}, error) {
	vls, err := d.vls()
	if err != nil {
		return nil, errors.Locus(err)
	}

	return c.Do("DEL", vls...)
}
