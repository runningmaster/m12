package core

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"

	"internal/errors"
)

// []string -> []interface{}
func valsFromStrings(v []string) []interface{} {
	res := make([]interface{}, 0, len(v))
	for i := range v {
		res = append(res, v[i])
	}
	return res
}

// JSON -> []string -> []interface{}
func valsFromJSONStrings(b []byte) ([]interface{}, error) {
	var v []string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return nil, errors.Locus(err)
	}
	return valsFromStrings(v), nil
}

// []int64 -> []interface{}
func valsFromInt64s(v []int64) []interface{} {
	res := make([]interface{}, 0, len(v))
	for i := range v {
		res = append(res, v[i])
	}
	return res
}

// JSON -> []int64 -> []interface{}
func valsFromJSONInt64s(b []byte) ([]interface{}, error) {
	var v []int64
	err := json.Unmarshal(b, &v)
	if err != nil {
		return nil, errors.Locus(err)
	}
	return valsFromInt64s(v), nil
}

// string -> [0]interface{}
func mergeKeyVals(key string, v ...interface{}) []interface{} {
	if key == "" {
		return v
	}

	res := make([]interface{}, 0, len(v)+1)
	return append(append(res, key), v...)
}

// cmd key -> mem1 mem2 mem3 mem4
func makeVector(cmd, key string, v ...interface{}) []interface{} {
	m := 2
	if key == "" {
		m--
	}
	res := make([]interface{}, 0, len(v)+m)
	return append(append(res, cmd), mergeKeyVals(key, v...)...)
}

// cmd key -> mem1
// cmd key -> mem2
// cmd key -> mem3
func makeBulkOne(cmd, key string, v ...interface{}) [][]interface{} {
	res := make([][]interface{}, 0, len(v))
	for i := range v {
		res = append(res, makeVector(cmd, key, v[i]))
	}
	return res
}

// cmd key1 -> [key1:] mem1 mem2 mem3 mem4
// cmd key2 -> [key2:] mem1 mem2 mem3 mem4
// cmd key3 -> [key3:] mem1 mem2 mem3 mem4
func makeBulkVector(src [][]interface{}, cmd, key string, v ...interface{}) [][]interface{} {
	return append(src, makeVector(cmd, key, v...))
}

// TODO: parse keyspace
//	"keyspace": {
//		"db0": "keys=1,expires=0,avg_ttl=0"
//	},
func parseInfo(b []byte) (map[string]map[string]string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(b))
	mapper := make(map[string]map[string]string)

	var (
		line  string
		sect  string
		split []string
	)
	for scanner.Scan() {
		line = strings.ToLower(scanner.Text())
		if line == "" {
			continue
		}
		if strings.Contains(line, "#") {
			sect = line[2:]
			mapper[sect] = make(map[string]string)
			continue
		}
		split = strings.Split(line, ":")
		mapper[sect][split[0]] = split[1]
	}

	if scanner.Err() != nil {
		return nil, errors.Locus(scanner.Err())
	}

	return mapper, nil
}
