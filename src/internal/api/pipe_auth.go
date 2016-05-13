package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"internal/conf"
	"internal/core"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
)

func pipeAuth(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			tkn *jwt.Token
			key string
			src []byte
			err error
		)

		tkn, err = getJWT(r)
		if err != nil {
			goto fail
		}

		key = tkn.Claims["skey"].(string)
		err = auth(key)
		if err != nil {
			goto fail
		}

		src, err = json.Marshal(tkn.Claims)
		if err != nil {
			goto fail
		}

		h(withAuth(withMeta(ctx, string(src)), key), w, r)
		return // success
	fail:
		h(withCode(withFail(ctx, err), http.StatusInternalServerError), w, r)
	}
}

func getJWT(r *http.Request) (*jwt.Token, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(conf.JWTSecretKey), nil
	}
	return jwt.ParseFromRequest(r, keyFunc)
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

	return fmt.Errorf("api: auth key not found: %s: forbidden", string(b))
}

func isMasterKey(key string) bool {
	return strings.EqualFold(conf.Masterkey, key)
}
