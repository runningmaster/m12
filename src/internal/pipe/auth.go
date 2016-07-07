package pipe

import (
	"fmt"
	"net/http"
	"strings"

	"internal/ctxutil"
	"internal/pref"
)

func Auth(master int) handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			key, err := getKey(r)
			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctxutil.WithFail(ctx, err, http.StatusForbidden)))
				return
			}

			err = auth(key, master)
			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctxutil.WithFail(ctx, err, http.StatusForbidden)))
				return
			}

			ctx = ctxutil.WithAuth(ctx, key)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
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

	ok, err := false, fmt.Errorf("FIXME") //api.AuthOK(key)
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
