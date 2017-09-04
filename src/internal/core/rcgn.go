package core

import (
	"encoding/json"
	"strings"
)

func Rcgn(meta, data []byte) (interface{}, error) {
	m, err := unmarshalMeta(meta)
	if err != nil {
		return nil, err
	}

	err = testHTag(m.HTag)
	if err != nil {
		return nil, err
	}

	v, err := unmarshalRcgn(data, m)
	if err != nil {
		return nil, err
	}

	d, err := mineLinks(v, m)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func isRcgnAddr(s string) bool {
	return strings.Contains(s, "rcgn.addr")
}

func unmarshalRcgn(data []byte, m *meta) (interface{}, error) {
	d := killUTF8BOM(data)
	m.ETag = btsToMD5(d)
	m.Size = int64(len(d))

	t := m.HTag

	switch {
	case isRcgnAddr(t):
		v := jsonRcgnAddr{}
		err := json.Unmarshal(data, &v)
		return v, err
	default:
		v := jsonRcgnDrug{}
		err := json.Unmarshal(data, &v)
		return v, err
	}
}
