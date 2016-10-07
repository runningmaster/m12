package pipe

import (
	"fmt"
	"net/http"
	"strings"

	"internal/core/ctxt"
)

func ErrH(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		code, text := codeWithText(r.URL.Path)
		ctx = ctxt.WithFail(ctx, fmt.Errorf("api: %s", text), code)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func codeWithText(s string) (int, string) {
	code := http.StatusInternalServerError
	switch {
	case strings.HasPrefix(s, "404"):
		code = http.StatusNotFound
	case strings.HasPrefix(s, "405"):
		code = http.StatusMethodNotAllowed
	}
	return code, strings.ToLower(http.StatusText(code))
}
