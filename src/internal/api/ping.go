package api

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/garyburd/redigo/redis"
)

// Ping calls Redis PING
func Ping(_ []byte) (interface{}, error) {
	c := redisConn()
	defer closeConn(c)

	return c.Do("PING")
}

// Info calls Redis INFO
func Info(_ []byte) (interface{}, error) {
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
