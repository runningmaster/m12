package server

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"internal/errors"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
)

type regHandler struct {
	m string
	p string
	h http.Handler
}

var regHandlers []regHandler

// RegHandler is called from another packages
func RegHandler(m, p string, h http.Handler) {
	regHandlers = append(regHandlers, regHandler{m, p, h})
}

// initRouter is work-around wrapper for router in *echo.Echo
func initRouter(r *echo.Echo, reg ...regHandler) error {
	for _, v := range reg {
		switch v.m {
		case echo.GET:
			r.Get(v.p, standard.WrapHandler(v.h))
		case echo.POST:
			r.Post(v.p, standard.WrapHandler(v.h))
		default:
			return errors.Locusf("router: unsupported method")
		}
	}
	return nil
}

// trapErrorHandler replaces echo.DefaultHTTPErrorHandler() with workaround for 404 and 405 errors
func trapErrorHandler(err error, c echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok && !c.Response().Committed() {
		switch he.Code {
		case http.StatusNotFound:
			c.Echo().Router().Find("GET", "/error/404", c)
			_ = c.Handle(c)
		case http.StatusMethodNotAllowed:
			c.Echo().Router().Find("GET", "/error/405", c)
			_ = c.Handle(c)
		}
		return
	}
	c.Echo().DefaultHTTPErrorHandler(err, c)
}

func makeRouter() (*echo.Echo, error) {
	r := echo.New()
	r.SetLogOutput(ioutil.Discard)
	r.SetHTTPErrorHandler(trapErrorHandler)

	err := initRouter(r, regHandlers...)
	if err != nil {
		return nil, errors.Locus(err)
	}

	return r, nil
}

func withRouter(addr string, r engine.Handler) (*http.Server, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, errors.Locus(err)
	}

	s := standard.New(u.Host)
	s.SetHandler(r)

	return s.Server, nil
}
