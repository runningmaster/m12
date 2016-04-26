package core

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// Redis scheme:
// HASH key => [l/v] [a/v] [s/v] [e/v]

type decodeLinkAddr []byte

// Redis scheme:
// HASH => key=ID (SHA1)
// HMSET key l/v a/v s/v e/v (if exists in json)
// HMGET key l a s e
// JSON array: [{"id":"key1","id_link":1,"id_addr":2,"id_stat":0,"egrpou":"egrpou1"}]
type linkAddr struct {
	ID     string `json:"id,omitempty"      redis:"key"`
	IDLink int64  `json:"id_link,omitempty" redis:"l"`
	IDAddr int64  `json:"id_addr,omitempty" redis:"a"`
	IDStat int64  `json:"id_stat,omitempty" redis:"s"`
	EGRPOU string `json:"egrpou,omitempty"  redis:"e"`
}

func (d decodeLinkAddr) src() ([]string, error) {
	var out []string
	err := json.Unmarshal(d, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (d decodeLinkAddr) lnk() ([]linkAddr, error) {
	var out []linkAddr
	err := json.Unmarshal(d, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (d decodeLinkAddr) vls() ([]interface{}, error) {
	src, err := d.src()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	var l linkAddr
	for i := range src {
		err = c.Send("HMGET", l.keyflds(src[i])...)
		if err != nil {
			return nil, err
		}
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	out := make([]interface{}, 0, len(src))
	var rcv []interface{}
	for i := range src {
		rcv, err = redis.Values(c.Receive())
		if err != nil {
			return nil, err
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
		err = c.Send("DEL", lnk[i].ID)
		if err != nil {
			return nil, err
		}
		err = c.Send("HMSET", lnk[i].keyvals()...)
		if err != nil {
			return nil, err
		}
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	for range lnk {
		_, err = c.Receive()
		if err != nil {
			return nil, err
		}
	}

	return stringOK(), nil
}

func (d decodeLinkAddr) del(c redis.Conn) (interface{}, error) {
	vls, err := d.vls()
	if err != nil {
		return nil, err
	}

	return c.Do("DEL", vls...)
}

func (l linkAddr) keyflds(k string) []interface{} {
	return []interface{}{k, "l", "a", "s", "e"}
}

func (l linkAddr) keyvals() []interface{} {
	v := make([]interface{}, 0, 1+4*2)
	v = append(v, l.ID)
	if l.IDLink != 0 {
		v = append(v, "l", l.IDLink)
	}
	if l.IDAddr != 0 {
		v = append(v, "a", l.IDAddr)
	}
	if l.IDStat != 0 {
		v = append(v, "s", l.IDStat)
	}
	if l.EGRPOU != "" {
		v = append(v, "e", l.EGRPOU)
	}
	return v
}

func (l linkAddr) makeFrom(k string, v []interface{}) interface{} {
	if isEmpty(v) {
		return nil
	}
	return linkAddr{
		ID:     k,
		IDLink: toInt64(v[0]),  // "l"
		IDAddr: toInt64(v[1]),  // "a"
		IDStat: toInt64(v[2]),  // "s"
		EGRPOU: toString(v[3]), // "e"
	}
}

func findLinkAddr(keys ...string) ([]linkAddr, error) {
	b, err := json.Marshal(keys)
	if err != nil {
		return nil, err
	}

	c := redisGet()
	defer redisPut(c)

	v, err := decodeLinkAddr(b).get(c)
	if err != nil {
		return nil, err
	}

	res := make([]linkAddr, 0, len(v))
	for i := range v {
		switch l := v[i].(type) {
		case nil:
			res = append(res, linkAddr{})
		case linkAddr:
			res = append(res, l)
		}
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("core: link addr: unreachable")
	}

	return res, nil
}
