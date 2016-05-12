package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"internal/core"
	"internal/flag"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
)

func pipeAuth(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		key, err := getKey(r)
		if err != nil {
			goto fail
		}

		err = auth(key)
		if err != nil {
			goto fail
		}

		h(withAuth(ctx, key), w, r)
		return // success
	fail:
		h(withCode(withFail(ctx, err), http.StatusInternalServerError), w, r)
	}
}

func getKeyV1(r *http.Request) (string, bool) {
	key := r.FormValue("key")
	return key, key != ""
}

func getKeyV2(r *http.Request) (string, bool) {
	key := r.Header.Get("X-Morion-Skynet-Key")
	return key, key != ""
}

// api:key-3ax6xnjp29jd6fds4gc373sgvjxteol0 (?)
func getKeyV3(r *http.Request) (string, bool) {
	var key string
	if _, pass, ok := r.BasicAuth(); ok {
		key = pass[4:]
	}

	return key, key != ""
}

// JWT experiment
func getKeyV4(r *http.Request) (string, bool) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(flag.JWTSecretKey), nil
	}
	t, err := jwt.ParseFromRequest(r, keyFunc)
	if err != nil {
		return "", false
	}

	if v, ok := t.Header["skey"].(string); ok {
		return v, true
	}

	return "", false
}

func getKey(r *http.Request) (string, error) {
	var (
		key string
		ok  bool
	)

	if key, ok = getKeyV4(r); ok {
		goto success
	}
	if key, ok = getKeyV3(r); ok {
		goto success
	}
	if key, ok = getKeyV2(r); ok {
		goto success
	}
	if key, ok = getKeyV1(r); ok {
		goto success
	}

	return "", fmt.Errorf("api: auth key (as param) not found")

success:
	return key, nil
}

func auth(key string) error {
	req, err := http.NewRequest("", "", strings.NewReader(fmt.Sprintf("[%q]", key)))
	if err != nil {
		return err
	}

	res, err := core.RunC("get", "auth")(context.Background(), nil, req)
	if err != nil {
		return err
	}

	if src, ok := res.([]interface{}); ok && len(src) > 0 {
		if val, ok := src[0].(string); ok && strings.EqualFold(val, key) {
			return nil
		}
	}

	if isMasterKey(key) {
		return nil
	}

	b, err := json.Marshal(res)
	if err != nil {
		return err
	}

	return fmt.Errorf("api: auth key (as value) not found: %s: forbidden", string(b))
}

func isMasterKey(key string) bool {
	return strings.EqualFold(flag.Masterkey, key)
}
