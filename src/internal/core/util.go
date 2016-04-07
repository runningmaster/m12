package core

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"

	"internal/compress/gzip"

	"github.com/spkg/bom"
)

func stringOK() string {
	return http.StatusText(http.StatusOK)
}

func mendGzip(ct string, b []byte) ([]byte, error) {
	if strings.Contains(ct, "gzip") {
		unz, err := gzip.Gunzip(b)
		if err != nil {
			return nil, err
		}
		return unz, nil
	}

	return b, nil
}

func mendUTF8(ct string, b []byte) ([]byte, error) {
	if strings.Contains(ct, "text/plain; charset=utf-8") {
		return bom.Clean(b), nil
	}

	return b, nil
}

func mendGzipAndUTF8(b []byte) ([]byte, error) {
	var (
		ct  string
		err error
	)

	ct = http.DetectContentType(b)
	if b, err = mendGzip(ct, b); err != nil {
		return nil, err
	}

	ct = http.DetectContentType(b)
	if b, err = mendUTF8(ct, b); err != nil {
		return nil, err
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
		return nil, scanner.Err()
	}

	return mapper, nil
}
