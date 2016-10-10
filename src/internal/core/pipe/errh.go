package pipe

import (
	"fmt"
	"net/http"
	"strings"
)

func ErrH(code int) handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = withFail(ctx, fmt.Errorf("pipe: %s", strings.ToLower(http.StatusText(code))), code)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
