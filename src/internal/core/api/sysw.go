package api

import (
	"encoding/json"
	"net/http"

	"internal/core/link"
	"internal/core/redis"
)

func getAuth(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return redis.GetLinkAuth(v)
}

func setAuth(data []byte) (interface{}, error) {
	var v []link.Auth
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return redis.SetLinkAuth(v)
}

func delAuth(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return redis.DelLinkAuth(v)
}

func getAddr(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return redis.GetLinkAddr(v)
}

func setAddr(data []byte) (interface{}, error) {
	var v []link.Addr
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return redis.SetLinkAddr(v)
}

func delAddr(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return redis.DelLinkAddr(v)
}

func getDrug(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return redis.GetLinkDrug(v)
}

func setDrug(data []byte) (interface{}, error) {
	var v []link.Drug
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return redis.SetLinkDrug(v)
}

func delDrug(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return redis.DelLinkDrug(v)
}

func getStat(data []byte, _, _ http.Header) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return redis.GetLinkStat(v)
}

func setStat(data []byte, _, _ http.Header) (interface{}, error) {
	var v []link.Stat
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return redis.SetLinkStat(v)
}

func delStat(data []byte, _, _ http.Header) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return redis.DelLinkStat(v)
}
