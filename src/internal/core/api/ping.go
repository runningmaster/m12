package api

import (
	"bufio"
	"bytes"
	"strings"

	"internal/database/redis"
)

func ping() (interface{}, error) {
	c := redis.Conn()
	defer redis.Free(c)

	return c.Do("PING")
}

func info() (interface{}, error) {
	c := redis.Conn()
	defer redis.Free(c)

	b, err := redis.Conv.ToBytes(c.Do("INFO"))
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
