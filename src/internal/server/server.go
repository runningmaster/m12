package server

import (
	"fmt"
	"internal/api"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/tylerb/graceful"
)

type regHandler struct {
	m string
	p string
	h http.Handler
}

var regHandlers []regHandler

// Run starts a server with router
func Run(addr string) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	r, err := makeRouter()
	if err != nil {
		return err
	}

	s, err := makeServer(u.Host, r)
	if err != nil {
		return err
	}

	return s.ListenAndServe()
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
			return fmt.Errorf("server: unsupported method")
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
			goto find
		case http.StatusMethodNotAllowed:
			c.Echo().Router().Find("GET", "/error/405", c)
			goto find
		}
	}
	c.Echo().DefaultHTTPErrorHandler(err, c)
	return
find:
	_ = c.Handler()(c)
}

func makeRouter() (*echo.Echo, error) {
	err := api.Reg(func(m, p string, h http.Handler) {
		regHandlers = append(regHandlers, regHandler{m, p, h})
	})
	if err != nil {
		return nil, err
	}

	r := echo.New()
	r.SetLogOutput(ioutil.Discard)
	r.SetHTTPErrorHandler(trapErrorHandler)

	err = initRouter(r, regHandlers...)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func makeServer(addr string, r *echo.Echo) (*graceful.Server, error) {
	s := standard.New(addr)
	s.SetHandler(r)

	return &graceful.Server{
			Server:  s.Server,
			Timeout: 5 * time.Second,
			Logger:  log.New(os.Stderr, "server: ", 0)},
		nil
}
