package core

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

// Redis scheme:
// HASH key => k/v [k/v ...]

const keyLinkStat = "link:stat"

type (
	decodeLinkStat []byte

	// Redis scheme:
	// HASH => key="stat"
	// HMSET key i->n [i->n...]
	// HMGET key i [i..]
	linkStat struct {
		ID   int64  `json:"id,omitempty"   redis:"i"`
		Name string `json:"name,omitempty" redis:"n"`
	}
)

func (d decodeLinkStat) src() ([]int64, error) {
	var out []int64
	err := json.Unmarshal(d, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (d decodeLinkStat) lnk() ([]linkStat, error) {
	var out []linkStat
	err := json.Unmarshal(d, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (d decodeLinkStat) vls(withKey bool) ([]interface{}, error) {
	src, err := d.src()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	rcv, err := redis.Values(c.Do("HMGET", vls...))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	vls := make([]interface{}, 0, len(lnk)*2+1)
	vls = append(vls, keyLinkStat)
	for i := range lnk {
		vls = append(vls, lnk[i].ID, lnk[i].Name)
	}

	_, err = c.Do("HMSET", vls...)
	if err != nil {
		return nil, err
	}

	return stringOK(), nil
}

func (d decodeLinkStat) del(c redis.Conn) (interface{}, error) {
	vls, err := d.vls(true)
	if err != nil {
		return nil, err
	}

	return c.Do("HDEL", vls...)
}

func (l linkStat) makeFrom(k int64, v interface{}) interface{} {
	if v == nil {
		return v
	}
	return linkStat{
		ID:   k,
		Name: toString(v), // "n"
	}
}
