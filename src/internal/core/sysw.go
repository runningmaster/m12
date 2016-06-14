package core

import (
	"encoding/json"
	"fmt"

	"internal/redis"
)

const (
	keyAuth = "list:auth"
	keyStat = "list:stat"
)

var (
	fldsAddr = [...]string{"l", "a", "s", "e"}
	fldsDrug = [...]string{"l", "d", "b", "c", "s"}
)

func AuthOK(key string) (bool, error) {
	v, err := redis.SISMEMBER(keyAuth, key)
	if err != nil {
		return false, err
	}

	return redis.ToInt64Safely(v) == 1, nil
}

// GetAuth returns
func GetAuth(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return getAuth(v...)
}

// SetAuth returns
func SetAuth(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return setAuth(v...)
}

// DelAuth returns
func DelAuth(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return delAuth(v...)
}

func getAuth(v ...string) ([]interface{}, error) {
	return redis.SISMEMBERM(redis.ConvFromStringsWithKey(keyAuth, v...)...)
}

func setAuth(v ...string) (interface{}, error) {
	return redis.SADD(redis.ConvFromStringsWithKey(keyAuth, v...)...)
}

func delAuth(v ...string) (interface{}, error) {
	return redis.SREM(redis.ConvFromStringsWithKey(keyAuth, v...)...)
}

// GetLinkAddr returns
func GetLinkAddr(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return getLinkAddr(v...)
}

// SetLinkAddr returns
func SetLinkAddr(data []byte) (interface{}, error) {
	var v []linkAddr
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return setLinkAddr(v...)
}

// DelLinkAddr returns
func DelLinkAddr(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return delLinkAddr(v...)
}

func valIsNill(v ...interface{}) bool {
	for i := range v {
		if v[i] != nil { // don't work with pointers!
			return false
		}
	}

	return true
}

func getLinkAddr(v ...string) ([]linkAddr, error) {
	vlm := make([][]interface{}, len(v))
	var vls []interface{}
	for i := range vlm {
		vls = make([]interface{}, len(fldsAddr)+1)
		vls[0] = v[i]        // key
		vls[1] = fldsAddr[0] // fld "l"
		vls[2] = fldsAddr[1] // fld "a"
		vls[3] = fldsAddr[2] // fld "s"
		vls[4] = fldsAddr[3] // fld "e"
		vlm[i] = vls
	}

	vlm, err := redis.HMGETM(vlm...)
	if err != nil {
		return nil, err
	}
	if len(vlm) != len(v) {
		return nil, fmt.Errorf("core: invalid len (get link addr): got %d, want %d", len(vlm), len(v))
	}

	res := make([]linkAddr, len(vlm))
	for i := range vlm {
		res[i].ID = v[i] // key
		if valIsNill(vlm[i]...) {
			continue
		}
		if len(vlm[i]) != len(fldsAddr) {
			return nil, fmt.Errorf("core: invalid len (get link addr): got %d, want %d", len(vlm[i]), len(fldsAddr))
		}
		res[i].IDLink = redis.ToInt64Safely(vlm[i][0])  // fld "l"
		res[i].IDAddr = redis.ToInt64Safely(vlm[i][1])  // fld "a"
		res[i].IDStat = redis.ToInt64Safely(vlm[i][2])  // fld "s"
		res[i].EGRPOU = redis.ToStringSafely(vlm[i][3]) // fld "e"
	}

	return res, nil
}

func setLinkAddr(v ...linkAddr) (interface{}, error) {
	vlm := make([][]interface{}, len(v))
	var vls []interface{}
	for i := range vlm {
		vls = make([]interface{}, 0, len(fldsAddr)*2+1)
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
		vlm[i] = vls
	}

	return redis.HMSETM(vlm...)
}

func delLinkAddr(v ...string) (interface{}, error) {
	return redis.DEL(redis.ConvFromStrings(v...)...)
}

// GetLinkDrug returns
func GetLinkDrug(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return getLinkDrug(v...)
}

// SetLinkDrug returns
func SetLinkDrug(data []byte) (interface{}, error) {
	var v []linkDrug
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return setLinkDrug(v...)
}

// DelLinkDrug returns
func DelLinkDrug(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return delLinkDrug(v...)
}

func getLinkDrug(v ...string) ([]linkDrug, error) {
	vlm := make([][]interface{}, len(v))
	var vls []interface{}
	for i := range vlm {
		vls = make([]interface{}, len(fldsDrug)+1)
		vls[0] = v[i]        // key
		vls[1] = fldsDrug[0] // fld "l"
		vls[2] = fldsDrug[1] // fld "d"
		vls[3] = fldsDrug[2] // fld "b"
		vls[4] = fldsDrug[3] // fld "c"
		vls[5] = fldsDrug[4] // fld "s"
		vlm[i] = vls
	}

	vlm, err := redis.HMGETM(vlm...)
	if err != nil {
		return nil, err
	}
	if len(vlm) != len(v) {
		return nil, fmt.Errorf("core: invalid len (get link drug): got %d, want %d", len(vlm), len(v))
	}

	res := make([]linkDrug, len(vlm))
	for i := range vlm {
		res[i].ID = v[i] // key
		if valIsNill(vlm[i]...) {
			continue
		}
		if len(vlm[i]) != len(fldsDrug) {
			return nil, fmt.Errorf("core: invalid len (get link drug): got %d, want %d", len(vlm[i]), len(fldsDrug))
		}
		res[i].IDLink = redis.ToInt64Safely(vlm[i][0]) // fld "l"
		res[i].IDDrug = redis.ToInt64Safely(vlm[i][1]) // fld "d"
		res[i].IDBrnd = redis.ToInt64Safely(vlm[i][2]) // fld "b"
		res[i].IDCatg = redis.ToInt64Safely(vlm[i][3]) // fld "c"
		res[i].IDStat = redis.ToInt64Safely(vlm[i][4]) // fld "s"
	}

	return res, nil
}

func setLinkDrug(v ...linkDrug) (interface{}, error) {
	vlm := make([][]interface{}, len(v))
	var vls []interface{}
	for i := range vlm {
		vls = make([]interface{}, 0, len(fldsDrug)*2+1)
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
		vlm[i] = vls
	}

	return redis.HMSETM(vlm...)
}

func delLinkDrug(v ...string) (interface{}, error) {
	return redis.DEL(redis.ConvFromStrings(v...)...)
}

// GetLinkStat returns
func GetLinkStat(data []byte) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return getLinkStat(v...)
}

// SetLinkStat returns
func SetLinkStat(data []byte) (interface{}, error) {
	var v []linkStat
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return setLinkStat(v...)
}

// DelLinkStat returns
func DelLinkStat(data []byte) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return delLinkStat(v...)
}

func getLinkStat(v ...int64) ([]linkStat, error) {
	vls, err := redis.HMGET(redis.ConvFromInt64sWithKey(keyStat, v...)...)
	if err != nil {
		return nil, err
	}

	if len(vls) != len(v) {
		return nil, fmt.Errorf("core: invalid len (get link stat): got %d, want %d", len(vls), len(v))
	}

	res := make([]linkStat, len(vls))
	for i := range vls {
		res[i].ID = v[i]
		if valIsNill(vls[i]) {
			continue
		}
		res[i].Name = redis.ToStringSafely(vls[i])
	}

	return res, nil
}

func setLinkStat(v ...linkStat) (interface{}, error) {
	vls := make([]interface{}, len(v)*2+1)
	vls[0] = keyStat
	for i := range v {
		vls[i*2+1] = v[i].ID
		vls[i*2+2] = v[i].Name
	}

	return redis.HMSET(vls...)
}

func delLinkStat(v ...int64) (interface{}, error) {
	return redis.HDEL(redis.ConvFromInt64sWithKey(keyStat, v...)...)
}
