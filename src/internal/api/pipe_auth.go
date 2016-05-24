package api

import (
	"fmt"
	"net/http"
	"strings"

	"internal/conf"
	"internal/core"

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

		h(ctxWithAuth(ctx, key), w, r)
		return // success
	fail:
		h(ctxWithCode(ctxWithFail(ctx, err), http.StatusForbidden), w, r)
	}
}

// api:key-3ax6xnjp29jd6fds4gc373sgvjxteol0 (?)
func getKey(r *http.Request) (string, error) {
	if _, pass, ok := r.BasicAuth(); ok {
		return pass[4:], nil
	}

	return "", fmt.Errorf("api: key not found")
}

func auth(key string) error {
	ok, err := core.AuthOK(key)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	if isMasterKey(key) {
		return nil
	}

	return fmt.Errorf("api: invalid key: %s: forbidden", key)
}

func isMasterKey(key string) bool {
	return strings.EqualFold(conf.Masterkey, key)
}
