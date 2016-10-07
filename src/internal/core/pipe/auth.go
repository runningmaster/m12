package pipe

import (
	"fmt"
	"net/http"
)

func Auth(fn func(string) bool) handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var ok bool
			key, err := getKey(r)
			if key != "" {
				ok = fn(key)
			}

			if ok {
				ctx = withAuth(ctx, key)
			} else {
				if err == nil {
					err = fmt.Errorf("pipe: invalid key: %s: forbidden", key)
				}
				ctx = withFail(ctx, err, http.StatusForbidden)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// V3 api:key-3ax6xnjp29jd6fds4gc373sgvjxteol0 (?)
func getKey(r *http.Request) (string, error) {
	if _, pass, ok := r.BasicAuth(); ok && len(pass) > 4 {
		return pass[4:], nil
	}
	return "", fmt.Errorf("pipe: key not found")
}
