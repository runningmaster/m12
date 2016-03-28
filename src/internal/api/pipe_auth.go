package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"internal/core"
	"internal/errors"
	"internal/flag"

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
		return

	fail:
		h(withCode(withFail(ctx, errors.Locus(err)), http.StatusInternalServerError), w, r)
	}
}

func getKeyV1(r *http.Request, err error) (string, error) {
	key := r.FormValue("key")
	if key == "" {
		return "", err
	}

	return key, nil
}

func getKeyV2(r *http.Request, err error) (string, error) {
	key := r.Header.Get("X-Morion-Skynet-Key")
	if key == "" {
		return "", err
	}

	return key, nil
}

func getKey(r *http.Request) (string, error) {
	var key string
	var err = errors.Locusf("api: auth key (as param) not found")

	if key, err = getKeyV1(r, err); err != nil {
		if key, err = getKeyV2(r, err); err != nil {
			return "", err
		}
	}

	return key, nil
}

func auth(key string) error {
	res, err := core.GetAuth([]byte(fmt.Sprintf("[%q]", key)))
	if err != nil {
		return errors.Locus(err)
	}

	if src, ok := res.([]interface{}); ok && len(src) > 0 {
		if val, ok := src[0].(string); ok && strings.EqualFold(val, key) {
			return nil
		}
	}

	if isMasterKey(key) {
		return nil
	}

	var v []byte
	if v, err = json.Marshal(res); err != nil {
		return errors.Locus(err)
	}
	return errors.Locusf("api: auth key (as value) not found: %s: forbidden", string(v))
}

func isMasterKey(key string) bool {
	return strings.EqualFold(flag.Masterkey, key)
}
