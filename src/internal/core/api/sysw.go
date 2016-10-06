package api

import (
	"encoding/json"
	"net/http"

	"internal/core"
)

func getAuth(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return core.GetLinkAuth(v)
}

func setAuth(data []byte) (interface{}, error) {
	var v []core.LinkAuth
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return core.SetLinkAuth(v)
}

func delAuth(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return core.DelLinkAuth(v)
}

func getAddr(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return core.GetLinkAddr(v)
}

func setAddr(data []byte) (interface{}, error) {
	var v []core.LinkAddr
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return core.SetLinkAddr(v)
}

func delAddr(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return core.DelLinkAddr(v)
}

func getDrug(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return core.GetLinkDrug(v)
}

func setDrug(data []byte) (interface{}, error) {
	var v []core.LinkDrug
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return core.SetLinkDrug(v)
}

func delDrug(data []byte) (interface{}, error) {
	var v []string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return core.DelLinkDrug(v)
}

func getStat(data []byte, _, _ http.Header) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return core.GetLinkStat(v)
}

func setStat(data []byte, _, _ http.Header) (interface{}, error) {
	var v []core.LinkStat
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return core.SetLinkStat(v)
}

func delStat(data []byte, _, _ http.Header) (interface{}, error) {
	var v []int64
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return core.DelLinkStat(v)
}
