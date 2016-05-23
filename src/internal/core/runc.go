package core

import (
	"encoding/json"
	"fmt"

	"internal/redis"
)

type opType int

const (
	opGet opType = iota
	opSet
	opDel

	keyAuth = "list:auth"
	keyStat = "list:stat"
)

func doAuth(data []byte, op opType) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	vls := redis.ConvFromStringsWithKey(keyAuth, v...)
	switch op {
	case opGet:
		return redis.SISMEMBERM(vls...)
	case opSet:
		return redis.SADD(vls...)
	case opDel:
		return redis.SREM(vls...)
	default:
		return nil, fmt.Errorf("core: unreachable")
	}
}

func doLinkAddr(data []byte, op opType) (interface{}, error) {
	var (
		v1   []*linkAddr
		v2   []string
		vls  []interface{}
		vlsm [][]interface{}
		flds = []string{"l", "a", "s", "e"}
		vals []interface{}
		err  error
	)

	if op == opSet {
		err = json.Unmarshal(data, &v1)
		if err != nil {
			return nil, err
		}
		vlsm = make([][]interface{}, len(v1))
		for i := range vlsm {
			if v1[i] == nil {
				continue
			}
			vals = make([]interface{}, len(flds)*2+1)
			vals[0] = v1[i].ID
			for j := range flds {
				vals[j*2+1] = flds[j]
				switch j {
				case 0:
					vals[j*2+2] = v1[i].IDLink
				case 1:
					vals[j*2+2] = v1[i].IDAddr
				case 2:
					vals[j*2+2] = v1[i].IDStat
				case 3:
					vals[j*2+2] = v1[i].EGRPOU
				}
			}
			vlsm[i] = vals
		}
	} else {
		err = json.Unmarshal(data, &v2)
		if err != nil {
			return nil, err
		}
		if op == opGet {
			vlsm = make([][]interface{}, len(v2))
			for i := range vlsm {
				vals = make([]interface{}, len(flds)+1)
				vals[0] = v2[i]
				for j := range flds {
					vals[j+1] = flds[j]
				}
			}
		} else {
			vls = redis.ConvFromStrings(v2...)
		}
	}

	switch op {
	case opGet:
		vlsm, err = redis.HMGETM(vlsm...)
		if err != nil {
			return nil, err
		}
		if len(vlsm) != len(v2) {
			return nil, fmt.Errorf("core: invalid len (get link addr): got %d, want %d", len(vlsm), len(v2))
		}
		res := make([]*linkAddr, len(vlsm))
		for i := range vlsm {
			if vlsm[i] == nil {
				continue
			}
			if len(vlsm[i]) != len(flds) {
				return nil, fmt.Errorf("core: invalid len (get link addr): got %d, want %d", len(vlsm[i]), len(flds))
			}
			res[i] = &linkAddr{
				ID: v2[i],
			}
			for j := range flds {
				switch j {
				case 0:
					res[i].IDLink = redis.ToInt64Safely(vlsm[i][j])
				case 1:
					res[i].IDAddr = redis.ToInt64Safely(vlsm[i][j])
				case 2:
					res[i].IDStat = redis.ToInt64Safely(vlsm[i][j])
				case 3:
					res[i].EGRPOU = redis.ToStringSafely(vlsm[i][j])
				}
			}
		}
		return res, nil
	case opSet:
		return redis.HMSETM(vlsm...)
	case opDel:
		return redis.DEL(vls...)
	}

	return nil, fmt.Errorf("core: unreachable")
}

func doLinkDrug(data []byte, op opType) (interface{}, error) {
	var (
		v1   []*linkDrug
		v2   []string
		vls  []interface{}
		vlsm [][]interface{}
		flds = []string{"l", "d", "b", "c", "s"}
		vals []interface{}
		err  error
	)

	if op == opSet {
		err = json.Unmarshal(data, &v1)
		if err != nil {
			return nil, err
		}
		vlsm = make([][]interface{}, len(v1))

		for i := range vlsm {
			if v1[i] == nil {
				continue
			}
			vals = make([]interface{}, len(flds)*2+1)
			vals[0] = v1[i].ID
			for j := range flds {
				vals[j*2+1] = flds[j]
				switch j {
				case 0:
					vals[j*2+2] = v1[i].IDLink
				case 1:
					vals[j*2+2] = v1[i].IDDrug
				case 2:
					vals[j*2+2] = v1[i].IDBrnd
				case 3:
					vals[j*2+2] = v1[i].IDCatg
				case 4:
					vals[j*2+2] = v1[i].IDStat
				}
			}
			vlsm[i] = vals
		}
	} else {
		err = json.Unmarshal(data, &v2)
		if err != nil {
			return nil, err
		}
		if op == opGet {
			vlsm = make([][]interface{}, len(v2))
			for i := range vlsm {
				vals = make([]interface{}, len(flds)+1)
				vals[0] = v2[i]
				for j := range flds {
					vals[j+1] = flds[j]
				}
			}
		} else {
			vls = redis.ConvFromStrings(v2...)
		}
	}

	switch op {
	case opGet:
		vlsm, err = redis.HMGETM(vlsm...)
		if err != nil {
			return nil, err
		}
		if len(vlsm) != len(v2) {
			return nil, fmt.Errorf("core: invalid len (get link drug): got %d, want %d", len(vlsm), len(v2))
		}
		res := make([]*linkDrug, len(vlsm))
		for i := range vlsm {
			if vlsm[i] == nil {
				continue
			}
			if len(vlsm[i]) != len(flds) {
				return nil, fmt.Errorf("core: invalid len (get link drug): got %d, want %d", len(vlsm[i]), len(flds))
			}
			res[i] = &linkDrug{
				ID: v2[i],
			}
			for j := range flds {
				switch j {
				case 0:
					res[i].IDLink = redis.ToInt64Safely(vlsm[i][j])
				case 1:
					res[i].IDDrug = redis.ToInt64Safely(vlsm[i][j])
				case 2:
					res[i].IDBrnd = redis.ToInt64Safely(vlsm[i][j])
				case 3:
					res[i].IDCatg = redis.ToInt64Safely(vlsm[i][j])
				case 4:
					res[i].IDStat = redis.ToInt64Safely(vlsm[i][j])
				}
			}
		}
		return res, nil
	case opSet:
		return redis.HMSETM(vlsm...)
	case opDel:
		return redis.DEL(vls...)
	}

	return nil, fmt.Errorf("core: unreachable")
}

func doLinkStat(data []byte, op opType) (interface{}, error) {
	var (
		v1  []*linkStat
		v2  []int64
		vls []interface{}
		err error
	)

	if op == opSet {
		err = json.Unmarshal(data, &v1)
		if err != nil {
			return nil, err
		}
		vls = make([]interface{}, len(v1)*2+1)
		vls[0] = keyStat
		for i := range v1 {
			vls[i*2+1] = v1[i].ID
			vls[i*2+2] = v1[i].Name
		}
	} else {
		err = json.Unmarshal(data, &v2)
		if err != nil {
			return nil, err
		}
		vls = redis.ConvFromInt64sWithKey(keyStat, v2...)
	}

	switch op {
	case opGet:
		vls, err = redis.HMGET(vls...)
		if err != nil {
			return nil, err
		}
		if len(vls) != len(v2) {
			return nil, fmt.Errorf("core: invalid len (get link stat): got %d, want %d", len(vls), len(v2))
		}
		res := make([]*linkStat, len(vls))
		for i := range vls {
			if vls[i] == nil {
				continue
			}
			res[i] = &linkStat{v2[i], redis.ToStringSafely(vls[i])}
		}
		return res, nil
	case opSet:
		return redis.HMSET(vls...)
	case opDel:
		return redis.HDEL(vls...)
	}

	return nil, fmt.Errorf("core: unreachable")
}

// GetAuth returns
func GetAuth(data []byte) (interface{}, error) {
	return doAuth(data, opGet)
}

// SetAuth returns
func SetAuth(data []byte) (interface{}, error) {
	return doAuth(data, opSet)
}

// DelAuth returns
func DelAuth(data []byte) (interface{}, error) {
	return doAuth(data, opDel)
}

// GetLinkAddr returns
func GetLinkAddr(data []byte) (interface{}, error) {
	return doLinkAddr(data, opGet)
}

// SetLinkAddr returns
func SetLinkAddr(data []byte) (interface{}, error) {
	return doLinkAddr(data, opSet)
}

// DelLinkAddr returns
func DelLinkAddr(data []byte) (interface{}, error) {
	return doLinkAddr(data, opDel)
}

// GetLinkDrug returns
func GetLinkDrug(data []byte) (interface{}, error) {
	return doLinkDrug(data, opGet)
}

// SetLinkDrug returns
func SetLinkDrug(data []byte) (interface{}, error) {
	return doLinkDrug(data, opSet)
}

// DelLinkDrug returns
func DelLinkDrug(data []byte) (interface{}, error) {
	return doLinkDrug(data, opDel)
}

// GetLinkStat returns
func GetLinkStat(data []byte) (interface{}, error) {
	return doLinkStat(data, opGet)
}

// SetLinkStat returns
func SetLinkStat(data []byte) (interface{}, error) {
	return doLinkStat(data, opSet)
}

// DelLinkStat returns
func DelLinkStat(data []byte) (interface{}, error) {
	return doLinkStat(data, opDel)
}

// RunC is "print-like" operation.
/*
func RunC(cmd, base string) Handler {
	return func(_ context.Context, _ http.ResponseWriter, r *http.Request) (interface{}, error) {
		b, err := readClose(r.Body)
		if err != nil {
			return nil, err
		}

		b, err = mendIfGzip(b)
		if err != nil {
			return nil, err
		}

		b, err = mendIfUTF8(b)
		if err != nil {
			return nil, err
		}

		v, err := makeGetSetDeler(base, b)
		if err != nil {
			return nil, err
		}

		return execGetSetDeler(cmd, v)
	}
}
*/
/*
func makeGetSetDeler(base string, b []byte) (redisGetSetDelOper, error) {
	switch base {
	case "auth":
		return nil decodeAuth(b), nil
	case "addr":
		return nil decodeLinkAddr(b), nil
	case "drug":
		return nil decodeLinkDrug(b), nil
	case "stat":
		return nil decodeLinkStat(b), nil
	}

	return nil, fmt.Errorf("core: unknown base %s", base)
}
*/

/*
func execGetSetDeler(cmd string, gsd redisGetSetDelOper) (interface{}, error) {
	c := redis2.Get()
	defer redis2.Put(c)

	switch cmd {
	case "get":
		return gsd.get(c)
	case "set":
		return gsd.set(c)
	case "del":
		return gsd.del(c)
	}

	return nil, fmt.Errorf("core: unknown command %s", cmd)
}
*/
