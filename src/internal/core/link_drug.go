package core

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// Redis scheme:
// HASH key => [l/v] [d/v] [b/v] [c/v] [s/v]

type decodeLinkDrug []byte

// Redis scheme:
// HASH => key=ID (SHA1)
// HMSET key l/v d/v b/v c/v s/v (if exists in json)
// HMGET key l d b c s
type linkDrug struct {
	ID     string `json:"id,omitempty"      redis:"key"`
	IDLink int64  `json:"id_link,omitempty" redis:"l"`
	IDDrug int64  `json:"id_drug,omitempty" redis:"d"`
	IDBrnd int64  `json:"id_brnd,omitempty" redis:"b"`
	IDCatg int64  `json:"id_catg,omitempty" redis:"c"`
	IDStat int64  `json:"id_stat,omitempty" redis:"s"`
}

func (d decodeLinkDrug) src() ([]string, error) {
	var out []string
	err := json.Unmarshal(d, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (d decodeLinkDrug) lnk() ([]linkDrug, error) {
	var out []linkDrug
	err := json.Unmarshal(d, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (d decodeLinkDrug) vls() ([]interface{}, error) {
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

func (d decodeLinkDrug) get(c redis.Conn) ([]interface{}, error) {
	src, err := d.src()
	if err != nil {
		return nil, err
	}

	var l linkDrug
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

func (d decodeLinkDrug) set(c redis.Conn) (interface{}, error) {
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

func (d decodeLinkDrug) del(c redis.Conn) (interface{}, error) {
	vls, err := d.vls()
	if err != nil {
		return nil, err
	}

	return c.Do("DEL", vls...)
}

func (l linkDrug) keyflds(k string) []interface{} {
	return []interface{}{k, "l", "d", "b", "c", "s"}
}

func (l linkDrug) keyvals() []interface{} {
	v := make([]interface{}, 0, 5*2)
	v = append(v, l.ID)
	if l.IDLink != 0 {
		v = append(v, "l", l.IDLink)
	}
	if l.IDDrug != 0 {
		v = append(v, "d", l.IDDrug)
	}
	if l.IDBrnd != 0 {
		v = append(v, "b", l.IDBrnd)
	}
	if l.IDCatg != 0 {
		v = append(v, "c", l.IDCatg)
	}
	if l.IDStat != 0 {
		v = append(v, "s", l.IDStat)
	}
	return v
}

func (l linkDrug) makeFrom(k string, v []interface{}) interface{} {
	if isEmpty(v) {
		return nil
	}
	return linkDrug{
		ID:     k,
		IDLink: toInt64(v[0]), // "l"
		IDDrug: toInt64(v[1]), // "d"
		IDBrnd: toInt64(v[2]), // "b"
		IDCatg: toInt64(v[3]), // "c"
		IDStat: toInt64(v[4]), // "s"
	}
}

func findLinkDrug(keys ...string) ([]linkDrug, error) {
	b, err := json.Marshal(keys)
	if err != nil {
		return nil, err
	}

	c := redisGet()
	defer redisPut(c)

	v, err := decodeLinkDrug(b).get(c)
	if err != nil {
		return nil, err
	}

	res := make([]linkDrug, 0, len(v))
	for i := range v {
		switch l := v[i].(type) {
		case nil:
			res = append(res, linkDrug{})
		case linkDrug:
			res = append(res, l)
		}
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("core: link drug: unreachable")
	}

	return res, nil
}
