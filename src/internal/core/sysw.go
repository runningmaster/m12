package core

import (
	"encoding/json"
	//"internal/redis"

	"github.com/garyburd/redigo/redis"
)

const (
	keyAuth = "list:auth"
	keyStat = "list:stat"
)

var (
	fldsAddr = []interface{}{"l", "a", "s", "e"}
	fldsDrug = []interface{}{"l", "d", "b", "c", "s"}
)

func valIsNill(v ...interface{}) bool {
	for i := range v {
		if v[i] != nil { // don't work with pointers!
			return false
		}
	}

	return true
}

func AuthOK(key string) (bool, error) {
	c := pREDIS.Get()
	defer c.Close()

	v, err := redis.Int(c.Do("HEXISTS", keyAuth, key))
	if err != nil {
		return false, err
	}

	return v == 1, nil
}

// GetAuth returns
func GetAuth(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return getAuth(v)
}

// SetAuth returns
func SetAuth(data []byte) (interface{}, error) {
	var v []linkAuth
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return setAuth(v)
}

// DelAuth returns
func DelAuth(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return delAuth(v)
}

func getAuth(v []string) ([]linkAuth, error) {
	c := pREDIS.Get()
	defer c.Close()

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
		if err != nil {
			return nil, err
		}
		out[i].Name = r
	}

	return out, nil
}

func setAuth(v []linkAuth) (interface{}, error) {
	c := pREDIS.Get()
	defer c.Close()

	var err error
	for i := range v {
		err = c.Send("HSET", keyAuth, v[i].ID, v[i].Name)
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func delAuth(v []string) (interface{}, error) {
	c := pREDIS.Get()
	defer c.Close()

	var err error
	for i := range v {
		err = c.Send("HDEL", keyAuth, v[i])
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

// GetAddr returns
func GetAddr(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return getAddr(v)
}

// SetAddr returns
func SetAddr(data []byte) (interface{}, error) {
	var v []linkAddr
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return setAddr(v)
}

// DelAddr returns
func DelAddr(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return delAddr(v)
}

func getAddr(v []string) ([]linkAddr, error) {
	c := pREDIS.Get()
	defer c.Close()

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
		if err != nil {
			return nil, err
		}
		if len(r) == len(fldsAddr) {
			out[i].IDLink, err = redis.Int64(r[0], err)  // fld "l"
			out[i].IDAddr, err = redis.Int64(r[1], err)  // fld "a"
			out[i].IDStat, err = redis.Int64(r[2], err)  // fld "s"
			out[i].EGRPOU, err = redis.String(r[3], err) // fld "e"
			if err != nil {
				return nil, err
			}
		}
	}

	return out, nil
}

func setAddr(v []linkAddr) (interface{}, error) {
	c := pREDIS.Get()
	defer c.Close()

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

func delAddr(v []string) (interface{}, error) {
	c := pREDIS.Get()
	defer c.Close()

	var err error
	for i := range v {
		err = c.Send("DEL", v[i])
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

// GetDrug returns
func GetDrug(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return getDrug(v)
}

// SetDrug returns
func SetDrug(data []byte) (interface{}, error) {
	var v []linkDrug
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return setDrug(v)
}

// DelDrug returns
func DelDrug(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return delDrug(v)
}

func getDrug(v []string) ([]linkDrug, error) {
	c := pREDIS.Get()
	defer c.Close()

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
		if err != nil {
			return nil, err
		}
		if len(r) == len(fldsDrug) {
			out[i].IDLink, err = redis.Int64(r[0], err) // fld "l"
			out[i].IDDrug, err = redis.Int64(r[1], err) // fld "d"
			out[i].IDBrnd, err = redis.Int64(r[2], err) // fld "b"
			out[i].IDCatg, err = redis.Int64(r[3], err) // fld "c"
			out[i].IDStat, err = redis.Int64(r[4], err) // fld "s"
			if err != nil {
				return nil, err
			}
		}
	}

	return out, nil

}

func setDrug(v []linkDrug) (interface{}, error) {
	c := pREDIS.Get()
	defer c.Close()

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

func delDrug(v []string) (interface{}, error) {
	c := pREDIS.Get()
	defer c.Close()

	var err error
	for i := range v {
		err = c.Send("DEL", v[i])
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

// GetStat returns
func GetStat(data []byte) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return getStat(v)
}

// SetStat returns
func SetStat(data []byte) (interface{}, error) {
	var v []linkStat
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return setStat(v)
}

// DelStat returns
func DelStat(data []byte) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return delStat(v)
}

func getStat(v []int64) ([]linkStat, error) {
	c := pREDIS.Get()
	defer c.Close()

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
		if err != nil {
			return nil, err
		}
		out[i].Name = r
	}

	return out, nil
}

func setStat(v []linkStat) (interface{}, error) {
	c := pREDIS.Get()
	defer c.Close()

	var err error
	for i := range v {
		err = c.Send("HSET", keyStat, v[i].ID, v[i].Name)
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}

func delStat(v []int64) (interface{}, error) {
	c := pREDIS.Get()
	defer c.Close()

	var err error
	for i := range v {
		err = c.Send("HDEL", keyStat, v[i])
		if err != nil {
			return nil, err
		}
	}

	return statusOK, c.Flush()
}
