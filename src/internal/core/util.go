package core

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"

	"internal/compress/gzip"
	"internal/errors"

	"github.com/spkg/bom"
)

func stringOK() string {
	return http.StatusText(http.StatusOK)
}

func mendGzip(b []byte) ([]byte, error) {
	if strings.Contains(http.DetectContentType(b), "gzip") {
		unz, err := gzip.Gunzip(b)
		if err != nil {
			return nil, errors.Locus(err)
		}
		return unz, nil
	}

	return b, nil
}

func mendBOM(b []byte) ([]byte, error) {
	if strings.Contains(http.DetectContentType(b), "text/plain; charset=utf-8") {
		return bom.Clean(b), nil
	}

	return b, nil
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
