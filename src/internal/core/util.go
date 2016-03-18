package core

import (
	"bufio"
	"bytes"
	"strings"

	"internal/errors"
)

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
