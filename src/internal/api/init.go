package api

import (
	"context"
	"net/http"
	"strings"

	"internal/core"

	"github.com/julienschmidt/httprouter"
)

// Init inits package vars
func Init() (http.Handler, error) {
	err := core.Init()
	if err != nil {
		return nil, err
	}
	return initRouter(), nil
}

func initRouter() http.Handler {
	r := httprouter.New()

	for k, v := range httpHandlers {
		s := strings.Split(k, ">") // [m,p]

		switch s[1] {
		case "/error/404":
			r.NotFound = v
		case "/error/405":
			r.MethodNotAllowed = v
		default:
			func(m, p string, h http.Handler) {
				r.Handle(m, p,
					func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
						ctx := r.Context()
						for i := range p {
							ctx = context.WithValue(ctx, p[i].Key, p[i].Value)
						}
						r = r.WithContext(ctx)
						h.ServeHTTP(w, r)
					})
			}(s[0], s[1], v)
		}
	}

	return r
}
