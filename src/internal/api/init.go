package api

import (
	"context"
	"net/http"
	"strings"

	"internal/core"

	"github.com/julienschmidt/httprouter"
)

func Init() (http.Handler, error) {
	err := core.Init()
	if err != nil {
		return nil, err
	}
	return initRouter(), nil
}

func initRouter() http.Handler {
	r := httprouter.New()
	httpWorkers = make(map[string]worker, len(httpHandlers))

	for k, v := range httpHandlers {
		s := strings.Split(k, ">") // [m,p]

		switch s[1] {
		case "/error/404":
			r.NotFound = v.h
		case "/error/405":
			r.MethodNotAllowed = v.h
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
			}(s[0], s[1], v.h)
		}

		if v.w != nil {
			httpWorkers[s[1]] = v.w
		}
	}

	return r
}
