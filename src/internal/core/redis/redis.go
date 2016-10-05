package redis

import (
	"encoding/json"
	"net/http"

	"internal/compress/gziputil"
	"internal/core/link"
	"internal/core/structs"
	"internal/database/redis"
)

const (
	keyAuth  = "list:auth" // FIXME: hset:auth
	keyStat  = "list:stat" // FIXME: hset:stat
	keyZlog  = "zset:meta"
	statusOK = http.StatusOK
)

var (
	fldsAddr = []interface{}{"l", "a", "s", "e"}
	fldsDrug = []interface{}{"l", "d", "b", "c", "s"}
)

func Pass(key string) bool {
	c := redis.Conn()
	defer redis.Free(c)

	v, _ := redis.Conv.ToInt64(c.Do("HEXISTS", keyAuth, key))

	return v == 1
}

func GetLinkAuth(v []string) ([]link.Auth, error) {
	c := redis.Conn()
	defer redis.Free(c)

	var err error
	for i := range v {
		err = c.Send("HGET", keyAuth, v[i])
		if err != nil {
			return nil, err
		}
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	out := make([]link.Auth, len(v))
	var r string
	for i := range v {
		out[i].ID = v[i]
		r, err = redis.Conv.ToString(c.Receive())
		if err != nil && redis.Conv.NotErrNil(err) {
			return nil, err
		}
		out[i].Name = r
	}

	return out, nil
}

func SetLinkAuth(v []link.Auth) (interface{}, error) {
	c := redis.Conn()
	defer redis.Free(c)

	var err error
	for i := range v {
		err = c.Send("HSET", keyAuth, v[i].ID, v[i].Name)
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func DelLinkAuth(v []string) (interface{}, error) {
	c := redis.Conn()
	defer redis.Free(c)

	var err error
	for i := range v {
		err = c.Send("HDEL", keyAuth, v[i])
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func GetLinkAddr(v []string) ([]link.Addr, error) {
	c := redis.Conn()
	defer redis.Free(c)

	vls := make([]interface{}, 0, len(fldsAddr)+1)
	var err error
	for i := range v {
		vls = append(vls, v[i]) // key
		vls = append(vls, fldsAddr...)

		err = c.Send("HMGET", vls...)
		if err != nil {
			return nil, err
		}

		vls = vls[:0]
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	out := make([]link.Addr, len(v))
	var r []interface{}
	for i := range v {
		out[i].ID = v[i] // key
		r, err = redis.Conv.ToIntfs(c.Receive())
		if err != nil && redis.Conv.NotErrNil(err) {
			return nil, err
		}
		if len(r) == len(fldsAddr) {
			out[i].IDLink, _ = redis.Conv.ToInt64(r[0], nil)  // fld "l"
			out[i].IDAddr, _ = redis.Conv.ToInt64(r[1], nil)  // fld "a"
			out[i].IDStat, _ = redis.Conv.ToInt64(r[2], nil)  // fld "s"
			out[i].EGRPOU, _ = redis.Conv.ToString(r[3], nil) // fld "e"
		}
	}

	return out, nil
}

func SetLinkAddr(v []link.Addr) (interface{}, error) {
	c := redis.Conn()
	defer redis.Free(c)

	vls := make([]interface{}, 0, len(fldsAddr)*2+1)
	var err error
	for i := range v {
		vls = append(vls, v[i].ID) // key
		if v[i].IDLink != 0 {
			vls = append(vls, fldsAddr[0], v[i].IDLink) // fld "l"
		}
		if v[i].IDAddr != 0 {
			vls = append(vls, fldsAddr[1], v[i].IDAddr) // fld "a"
		}
		if v[i].IDStat != 0 {
			vls = append(vls, fldsAddr[2], v[i].IDStat) // fld "s"
		}
		if v[i].EGRPOU != "" {
			vls = append(vls, fldsAddr[3], v[i].EGRPOU) // fld "e"
		}

		err = c.Send("DEL", v[i].ID)
		if err != nil {
			return nil, err
		}

		err = c.Send("HMSET", vls...)
		if err != nil {
			return nil, err
		}

		vls = vls[:0]
	}

	return statusOK, c.Flush()
}

func DelLinkAddr(v []string) (interface{}, error) {
	c := redis.Conn()
	defer redis.Free(c)

	var err error
	for i := range v {
		err = c.Send("DEL", v[i])
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func GetLinkDrug(v []string) ([]link.Drug, error) {
	c := redis.Conn()
	defer redis.Free(c)

	vls := make([]interface{}, 0, len(fldsDrug)+1)
	var err error
	for i := range v {
		vls = append(vls, v[i]) // key
		vls = append(vls, fldsDrug...)

		err = c.Send("HMGET", vls...)
		if err != nil {
			return nil, err
		}

		vls = vls[:0]
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	out := make([]link.Drug, len(v))
	var r []interface{}
	for i := range v {
		out[i].ID = v[i] // key
		r, err = redis.Conv.ToIntfs(c.Receive())
		if err != nil && redis.Conv.NotErrNil(err) {
			return nil, err
		}
		if len(r) == len(fldsDrug) {
			out[i].IDLink, _ = redis.Conv.ToInt64(r[0], nil) // fld "l"
			out[i].IDDrug, _ = redis.Conv.ToInt64(r[1], nil) // fld "d"
			out[i].IDBrnd, _ = redis.Conv.ToInt64(r[2], nil) // fld "b"
			out[i].IDCatg, _ = redis.Conv.ToInt64(r[3], nil) // fld "c"
			out[i].IDStat, _ = redis.Conv.ToInt64(r[4], nil) // fld "s"
		}
	}

	return out, nil

}

func SetLinkDrug(v []link.Drug) (interface{}, error) {
	c := redis.Conn()
	defer redis.Free(c)

	vls := make([]interface{}, 0, len(fldsDrug)*2+1)
	var err error
	for i := range v {
		vls = append(vls, v[i].ID) // key
		if v[i].IDLink != 0 {
			vls = append(vls, fldsDrug[0], v[i].IDLink) // fld "l"
		}
		if v[i].IDDrug != 0 {
			vls = append(vls, fldsDrug[1], v[i].IDDrug) // fld "d"
		}
		if v[i].IDBrnd != 0 {
			vls = append(vls, fldsDrug[2], v[i].IDBrnd) // fld "b"
		}
		if v[i].IDCatg != 0 {
			vls = append(vls, fldsDrug[3], v[i].IDCatg) // fld "c"
		}
		if v[i].IDStat != 0 {
			vls = append(vls, fldsDrug[4], v[i].IDStat) // fld "s"
		}

		err = c.Send("DEL", v[i].ID)
		if err != nil {
			return nil, err
		}

		err = c.Send("HMSET", vls...)
		if err != nil {
			return nil, err
		}

		vls = vls[:0]
	}

	return statusOK, c.Flush()
}

func DelLinkDrug(v []string) (interface{}, error) {
	c := redis.Conn()
	defer redis.Free(c)

	var err error
	for i := range v {
		err = c.Send("DEL", v[i])
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func GetLinkStat(v []int64) ([]link.Stat, error) {
	c := redis.Conn()
	defer redis.Free(c)

	var err error
	for i := range v {
		err = c.Send("HGET", keyStat, v[i])
		if err != nil {
			return nil, err
		}
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	out := make([]link.Stat, len(v))
	var r string
	for i := range v {
		out[i].ID = v[i]
		r, err = redis.Conv.ToString(c.Receive())
		if err != nil && redis.Conv.NotErrNil(err) {
			return nil, err
		}
		out[i].Name = r
	}

	return out, nil
}

func SetLinkStat(v []link.Stat) (interface{}, error) {
	c := redis.Conn()
	defer redis.Free(c)

	var err error
	for i := range v {
		err = c.Send("HSET", keyStat, v[i].ID, v[i].Name)
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func DelLinkStat(v []int64) (interface{}, error) {
	c := redis.Conn()
	defer redis.Free(c)

	var err error
	for i := range v {
		err = c.Send("HDEL", keyStat, v[i])
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func SetZlog(m structs.Meta) error {
	c := redis.Conn()
	defer redis.Free(c)

	z, err := gziputil.Compress(m.Marshal())
	if err != nil {
		return err
	}

	err = c.Send("ZADD", keyZlog, m.Unix, m.UUID)
	if err != nil {
		return err
	}

	err = c.Send("SET", m.UUID, z, "EX", 60*60*24*3)
	if err != nil {
		return err
	}

	return c.Flush()
}

func GetZlog(data []byte) (interface{}, error) {
	c := redis.Conn()
	defer redis.Free(c)

	res, err := redis.Conv.ToStrings(c.Do("ZRANGEBYSCORE", keyZlog, "-inf", "+inf"))
	if err != nil {
		return nil, err
	}
	for i := range res {
		err = c.Send("GET", res[i])
		if err != nil {
			return nil, err
		}
	}

	err = c.Flush()
	if err != nil {
		return nil, err
	}

	out := make([]structs.Meta, 0, len(res))
	var r []byte
	var z []byte
	var m structs.Meta
	for range res {
		z, err = redis.Conv.ToBytes(c.Receive())
		if err != nil && redis.Conv.NotErrNil(err) {
			return nil, err
		}

		r, err = gziputil.Uncompress(z)
		if err != nil {
			return nil, err
		}

		m, err = structs.UnmarshalMeta(r)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}

	return out, nil
}

func GetMeta(data []byte) (interface{}, error) {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	c := redis.Conn()
	defer redis.Free(c)

	z, err := redis.Conv.ToBytes(c.Do("GET", v))
	if err != nil /*&& err != redis.ErrNil*/ {
		return nil, err
	}

	r, err := gziputil.Uncompress(z)
	if err != nil {
		return nil, err
	}

	m, err := structs.UnmarshalMeta(r)
	if err != nil {
		return nil, err
	}

	return m, nil
}
