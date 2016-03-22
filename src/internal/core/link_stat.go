package core

import (
	"encoding/json"

	"internal/errors"

	"github.com/garyburd/redigo/redis"
)

// Redis scheme:
// HASH key => k/v [k/v ...]

const keyLinkStat = "link:stat"

type decodeLinkStat []byte

func (d decodeLinkStat) src() ([]int64, error) {
	var out []int64
	if err := json.Unmarshal(d, &out); err != nil {
		return nil, errors.Locus(err)
	}

	return out, nil
}

func (d decodeLinkStat) lnk() ([]linkStat, error) {
	var out []linkStat
	if err := json.Unmarshal(d, &out); err != nil {
		return nil, errors.Locus(err)
	}

	return out, nil
}

func (d decodeLinkStat) vls(withKey bool) ([]interface{}, error) {
	src, err := d.src()
	if err != nil {
		return nil, errors.Locus(err)
	}

	out := make([]interface{}, 0, len(src)+1)
	if withKey {
		out = append(out, keyLinkStat)
	}

	for i := range src {
		out = append(out, src[i])
	}

	return out, nil
}

func (d decodeLinkStat) get(c redis.Conn) ([]interface{}, error) {
	vls, err := d.vls(true)
	if err != nil {
		return nil, errors.Locus(err)
	}

	rcv, err := redis.Values(c.Do("HMGET", vls...))
	if err != nil {
		return nil, errors.Locus(err)
	}

	var l linkStat
	out := make([]interface{}, 0, len(rcv))
	for i := range rcv {
		out = append(out, l.makeFrom(toInt64(vls[i+1]), rcv[i]))
	}

	return out, nil
}

func (d decodeLinkStat) set(c redis.Conn) (interface{}, error) {
	lnk, err := d.lnk()
	if err != nil {
		return nil, errors.Locus(err)
	}

	vls := make([]interface{}, 0, len(lnk)*2+1)
	vls = append(vls, keyLinkStat)
	for i := range lnk {
		vls = append(vls, lnk[i].ID, lnk[i].Name)
	}

	if _, err = c.Do("HMSET", vls...); err != nil {
		return nil, errors.Locus(err)
	}

	return stringOK(), nil
}

func (d decodeLinkStat) del(c redis.Conn) (interface{}, error) {
	vls, err := d.vls(true)
	if err != nil {
		return nil, errors.Locus(err)
	}

	return c.Do("HDEL", vls...)
}
