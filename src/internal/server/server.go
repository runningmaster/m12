package server

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"internal/api"

	"github.com/julienschmidt/httprouter"
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

	_, p, err := net.SplitHostPort(u.Host)
	if err != nil {
		return err
	}
	log.Printf("server: started and listening to :%s", p)

	return s.ListenAndServe()
}

func makeRouter() (http.Handler, error) {
	err := api.Reg(func(m, p string, h http.Handler) {
		regHandlers = append(regHandlers, regHandler{m, p, h})
	})
	if err != nil {
		return nil, err
	}

	r := httprouter.New()
	for _, v := range regHandlers {
		switch v.p {
		case "/error/404":
			r.NotFound = v.h
		case "/error/405":
			r.MethodNotAllowed = v.h
		default:
			r.Handler(v.m, v.p, v.h)
		}
	}

	return r, nil
}

func makeServer(addr string, h http.Handler) (*graceful.Server, error) {
	return &graceful.Server{
			Server: &http.Server{
				Addr:    addr,
				Handler: h,
			},
			Timeout: 5 * time.Second,
			Logger:  log.New(os.Stderr, "server: ", 0)},
		nil
}
