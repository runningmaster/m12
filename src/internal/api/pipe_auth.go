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

			key, err := getKey(r)
			if err != nil {
				h(w, r.WithContext(ctxWithCode(ctxWithFail(ctx, err), http.StatusForbidden)))
				return
			}

			err = auth(key, master)
			if err != nil {
				h(w, r.WithContext(ctxWithCode(ctxWithFail(ctx, err), http.StatusForbidden)))
				return
			}

			ctx = ctxWithAuth(ctx, key)
			h(w, r.WithContext(ctx))
		}
	}
}

// V3 api:key-3ax6xnjp29jd6fds4gc373sgvjxteol0 (?)
func getKey(r *http.Request) (string, error) {
	if _, pass, ok := r.BasicAuth(); ok {
		return pass[4:], nil
	}
	return "", fmt.Errorf("api: key not found")
}

func auth(key string, master int) error {
	if isMasterKey(key) {
		return nil
	}

	if master == 1 {
		return fmt.Errorf("api: must be master key: forbidden")
	}

	ok, err := core.AuthOK(key)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	return fmt.Errorf("api: invalid key: %s: forbidden", key)
}

func isMasterKey(key string) bool {
	return strings.EqualFold(pref.MasterKey, key)
}
