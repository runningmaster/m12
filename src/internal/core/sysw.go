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
	var v []*linkStat
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

func getLinkStat(v ...int64) ([]*linkStat, error) {
	vls, err := redis.HMGET(redis.ConvFromInt64sWithKey(keyStat, v...)...)
	if err != nil {
		return nil, err
	}

	if len(vls) != len(v) {
		return nil, fmt.Errorf("core: invalid len (get link stat): got %d, want %d", len(vls), len(v))
	}

	res := make([]*linkStat, len(vls))
	for i := range vls {
		if vls[i] == nil {
			continue
		}
		res[i] = &linkStat{v[i], redis.ToStringSafely(vls[i])}
	}

	return res, nil
}

func setLinkStat(v ...*linkStat) (interface{}, error) {
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
	var v []*linkAddr
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

func getLinkAddr(v ...string) ([]*linkAddr, error) {
	vlm := make([][]interface{}, len(v))
	var vls []interface{}
	for i := range vlm {
		vls = make([]interface{}, len(fldsAddr)+1)
		vls[0] = v[i]
		vls[1] = fldsAddr[0]
		vls[2] = fldsAddr[1]
		vls[3] = fldsAddr[2]
		vls[4] = fldsAddr[3]
		vlm[i] = vls
	}

	vlm, err := redis.HMGETM(vlm...)
	if err != nil {
		return nil, err
	}
	if len(vlm) != len(v) {
		return nil, fmt.Errorf("core: invalid len (get link addr): got %d, want %d", len(vlm), len(v))
	}

	res := make([]*linkAddr, len(vlm))
	for i := range vlm {
		if vlm[i] == nil {
			continue
		}
		if len(vlm[i]) != len(fldsAddr) {
			return nil, fmt.Errorf("core: invalid len (get link addr): got %d, want %d", len(vlm[i]), len(fldsAddr))
		}
		res[i] = &linkAddr{
			ID:     v[i],
			IDLink: redis.ToInt64Safely(vlm[i][0]),
			IDAddr: redis.ToInt64Safely(vlm[i][1]),
			IDStat: redis.ToInt64Safely(vlm[i][2]),
			EGRPOU: redis.ToStringSafely(vlm[i][3]),
		}
	}

	return res, nil
}

func setLinkAddr(v ...*linkAddr) (interface{}, error) {
	vlm := make([][]interface{}, len(v))
	var vls []interface{}
	for i := range vlm {
		if v[i] == nil {
			continue
		}
		vls = make([]interface{}, len(fldsAddr)*2+1)
		vls[0] = v[i].ID
		vls[1] = fldsAddr[0]
		vls[2] = v[i].IDLink
		vls[3] = fldsAddr[1]
		vls[4] = v[i].IDAddr
		vls[5] = fldsAddr[2]
		vls[6] = v[i].IDStat
		vls[7] = fldsAddr[3]
		vls[8] = v[i].EGRPOU
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
	var v []*linkDrug
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

func getLinkDrug(v ...string) ([]*linkDrug, error) {
	vlm := make([][]interface{}, len(v))
	var vls []interface{}
	for i := range vlm {
		vls = make([]interface{}, len(fldsDrug)+1)
		vls[0] = v[i]
		vls[1] = fldsDrug[0]
		vls[2] = fldsDrug[1]
		vls[3] = fldsDrug[2]
		vls[4] = fldsDrug[3]
		vls[5] = fldsDrug[4]
		vlm[i] = vls
	}

	vlm, err := redis.HMGETM(vlm...)
	if err != nil {
		return nil, err
	}
	if len(vlm) != len(v) {
		return nil, fmt.Errorf("core: invalid len (get link drug): got %d, want %d", len(vlm), len(v))
	}

	res := make([]*linkDrug, len(vlm))
	for i := range vlm {
		if vlm[i] == nil {
			continue
		}
		if len(vlm[i]) != len(fldsDrug) {
			return nil, fmt.Errorf("core: invalid len (get link drug): got %d, want %d", len(vlm[i]), len(fldsDrug))
		}
		res[i] = &linkDrug{
			ID:     v[i],
			IDLink: redis.ToInt64Safely(vlm[i][0]),
			IDDrug: redis.ToInt64Safely(vlm[i][1]),
			IDBrnd: redis.ToInt64Safely(vlm[i][2]),
			IDCatg: redis.ToInt64Safely(vlm[i][3]),
			IDStat: redis.ToInt64Safely(vlm[i][4]),
		}
	}

	return res, nil
}

func setLinkDrug(v ...*linkDrug) (interface{}, error) {
	vlm := make([][]interface{}, len(v))
	var vls []interface{}
	for i := range vlm {
		if v[i] == nil {
			continue
		}
		vls = make([]interface{}, len(fldsDrug)*2+1)
		vls[0] = v[i].ID
		vls[1] = fldsDrug[0]
		vls[2] = v[i].IDLink
		vls[3] = fldsDrug[1]
		vls[4] = v[i].IDDrug
		vls[5] = fldsDrug[2]
		vls[6] = v[i].IDBrnd
		vls[7] = fldsDrug[3]
		vls[8] = v[i].IDCatg
		vls[9] = fldsDrug[4]
		vls[10] = v[i].IDStat
		vlm[i] = vls
	}

	return redis.HMSETM(vlm...)
}

func delLinkDrug(v ...string) (interface{}, error) {
	return redis.DEL(redis.ConvFromStrings(v...)...)
}
