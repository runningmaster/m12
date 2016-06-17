package api

import (
	"fmt"
	"net/http"
	"strings"

	"internal/core"
	"internal/pref"

	"golang.org/x/net/context"
)

func pipeAuth(master int) handlerPipe {
	return func(h handlerFunc) handlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			key, err := getKey(r)
			if err != nil {
				goto fail
			}

			err = auth(key, master)
			if err != nil {
				goto fail
			}

			h(ctxWithAuth(ctx, key), w, r)
			return // success
		fail:
			h(ctxWithCode(ctxWithFail(ctxWithAuth(ctx, key), err), http.StatusForbidden), w, r)
		}
	}
}

func getKey(r *http.Request) (string, error) {
	// V3 api:key-3ax6xnjp29jd6fds4gc373sgvjxteol0 (?)
	if _, pass, ok := r.BasicAuth(); ok {
		return pass[4:], nil
	}

	return "", fmt.Errorf("api: key not found")
}

func auth(key string, master int) error {
	if isMasterKey(key) {
		return nil
	}

	err403 := fmt.Errorf("api: invalid key: %s: forbidden", key)
	if master == 1 {
		return err403
	}

	ok, err := core.AuthOK(key)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	return err403
}

func isMasterKey(key string) bool {
	return strings.EqualFold(pref.MasterKey, key)
}
