package api

import (
	"fmt"
	"net/http"
	"strings"

	"internal/core"
	"internal/pref"
)

func pipeAuth(master int) handlerPipe {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			key, err := auth(r, master)
			if err != nil {
				ctx = ctxWithFail(ctx, err)
				ctx = ctxWithCode(ctx, http.StatusForbidden)
			}
			ctx = ctxWithAuth(ctx, key)
			h(w, r.WithContext(ctx))
		}
	}
}

func auth(r *http.Request, master int) (string, error) {
	key, err := getKey(r)
	if err != nil {
		return key, err
	}

	if isMasterKey(key) {
		return key, nil
	}

	err403 := fmt.Errorf("api: invalid key: %s: forbidden", key)
	if master == 1 {
		return key, err403
	}

	ok, err := core.AuthOK(key)
	if err != nil {
		return key, err
	}

	if ok {
		return key, nil
	}

	return key, err403
}

// V3 api:key-3ax6xnjp29jd6fds4gc373sgvjxteol0 (?)
func getKey(r *http.Request) (string, error) {
	if _, pass, ok := r.BasicAuth(); ok {
		return pass[4:], nil
	}
	return "", fmt.Errorf("api: key not found")
}

func isMasterKey(key string) bool {
	return strings.EqualFold(pref.MasterKey, key)
}
