package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"internal/pref"

	"github.com/garyburd/redigo/redis"
)

const (
	keyAuth  = "list:auth" // FIXME: hset:auth
	keyStat  = "list:stat" // FIXME: hset:stat
	statusOK = http.StatusOK
)

var (
	fldsAddr = []interface{}{"l", "a", "s", "e"}
	fldsDrug = []interface{}{"l", "d", "b", "c", "s"}
)

func pass(key string) bool {
	if strings.EqualFold(pref.MasterKey, key) {
		return true
	}

	c := redisConn()
	defer closeConn(c)

	v, _ := redis.Int(c.Do("HEXISTS", keyAuth, key))

	return v == 1
}

func getAuth(data []byte, _, _ http.Header) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return getAuthREDIS(v)
}

func setAuth(data []byte, _, _ http.Header) (interface{}, error) {
	var v []linkAuth
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return setAuthREDIS(v)
}

func delAuth(data []byte, _, _ http.Header) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return delAuthREDIS(v)
}

func getAuthREDIS(v []string) ([]linkAuth, error) {
	c := redisConn()
	defer closeConn(c)

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

	out := make([]linkAuth, len(v))
	var r string
	for i := range v {
		out[i].ID = v[i]
		r, err = redis.String(c.Receive())
		if err != nil && err != redis.ErrNil {
			return nil, err
		}
		out[i].Name = r
	}

	return out, nil
}

func setAuthREDIS(v []linkAuth) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

	var err error
	for i := range v {
		err = c.Send("HSET", keyAuth, v[i].ID, v[i].Name)
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func delAuthREDIS(v []string) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

	var err error
	for i := range v {
		err = c.Send("HDEL", keyAuth, v[i])
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func getAddr(data []byte, _, _ http.Header) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return getAddrREDIS(v)
}

func setAddr(data []byte, _, _ http.Header) (interface{}, error) {
	var v []linkAddr
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return setAddrREDIS(v)
}

func delAddr(data []byte, _, _ http.Header) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return delAddrREDIS(v)
}

func getAddrREDIS(v []string) ([]linkAddr, error) {
	c := redisConn()
	defer closeConn(c)

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

	out := make([]linkAddr, len(v))
	var r []interface{}
	for i := range v {
		out[i].ID = v[i] // key
		r, err = redis.Values(c.Receive())
		if err != nil && err != redis.ErrNil {
			return nil, err
		}
		if len(r) == len(fldsAddr) {
			out[i].IDLink, _ = redis.Int64(r[0], nil)  // fld "l"
			out[i].IDAddr, _ = redis.Int64(r[1], nil)  // fld "a"
			out[i].IDStat, _ = redis.Int64(r[2], nil)  // fld "s"
			out[i].EGRPOU, _ = redis.String(r[3], nil) // fld "e"
		}
	}

	return out, nil
}

func setAddrREDIS(v []linkAddr) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

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

func delAddrREDIS(v []string) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

	var err error
	for i := range v {
		err = c.Send("DEL", v[i])
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func getDrug(data []byte, _, _ http.Header) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return getDrugREDIS(v)
}

func setDrug(data []byte, _, _ http.Header) (interface{}, error) {
	var v []linkDrug
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return setDrugREDIS(v)
}

func delDrug(data []byte, _, _ http.Header) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return delDrugREDIS(v)
}

func getDrugREDIS(v []string) ([]linkDrug, error) {
	c := redisConn()
	defer closeConn(c)

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

	out := make([]linkDrug, len(v))
	var r []interface{}
	for i := range v {
		out[i].ID = v[i] // key
		r, err = redis.Values(c.Receive())
		if err != nil && err != redis.ErrNil {
			return nil, err
		}
		if len(r) == len(fldsDrug) {
			out[i].IDLink, _ = redis.Int64(r[0], nil) // fld "l"
			out[i].IDDrug, _ = redis.Int64(r[1], nil) // fld "d"
			out[i].IDBrnd, _ = redis.Int64(r[2], nil) // fld "b"
			out[i].IDCatg, _ = redis.Int64(r[3], nil) // fld "c"
			out[i].IDStat, _ = redis.Int64(r[4], nil) // fld "s"
		}
	}

	return out, nil

}

func setDrugREDIS(v []linkDrug) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

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

func delDrugREDIS(v []string) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

	var err error
	for i := range v {
		err = c.Send("DEL", v[i])
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func getStat(data []byte, _, _ http.Header) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return getStatREDIS(v)
}

func setStat(data []byte, _, _ http.Header) (interface{}, error) {
	var v []linkStat
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return setStatREDIS(v)
}

func delStat(data []byte, _, _ http.Header) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return delStatREDIS(v)
}

func getStatREDIS(v []int64) ([]linkStat, error) {
	c := redisConn()
	defer closeConn(c)

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

	out := make([]linkStat, len(v))
	var r string
	for i := range v {
		out[i].ID = v[i]
		r, err = redis.String(c.Receive())
		if err != nil && err != redis.ErrNil {
			return nil, err
		}
		out[i].Name = r
	}

	return out, nil
}

func setStatREDIS(v []linkStat) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

	var err error
	for i := range v {
		err = c.Send("HSET", keyStat, v[i].ID, v[i].Name)
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func delStatREDIS(v []int64) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

	var err error
	for i := range v {
		err = c.Send("HDEL", keyStat, v[i])
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}
