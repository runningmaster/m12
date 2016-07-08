package api

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"internal/version"

	"github.com/garyburd/redigo/redis"
)

func root(_ []byte, _, _ http.Header) (interface{}, error) {
	return fmt.Sprintf("%s %s", version.AppName(), version.WithBuildInfo()), nil
}

func ping(_ []byte, _, _ http.Header) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

	return c.Do("PING")
}

func info(_ []byte, _, _ http.Header) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

	b, err := redis.Bytes(c.Do("INFO"))
	if err != nil {
		return nil, err
	}

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
