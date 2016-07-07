package pipe

import (
	"fmt"
	"net/http"

	"internal/ctxutil"
)

func Auth(authFunc ...func(string) (bool, error)) handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			key, err := getKey(r)
			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctxutil.WithFail(ctx, err, http.StatusForbidden)))
				return
			}

			var ok bool
			for i := range authFunc {
				ok, err = authFunc[i](key)
				if err != nil {
					next.ServeHTTP(w, r.WithContext(ctxutil.WithFail(ctx, err, http.StatusForbidden)))
					return
				}
				if ok {
					break
				}
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
	return "", fmt.Errorf("pipe: key not found")
}
