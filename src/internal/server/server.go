package server

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"internal/api"

	"github.com/tylerb/graceful"
)

// Run starts a server with router
func Run(addr string) error {
	h, err := api.Init()
	if err != nil {
		return err
	}

	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	return makeServer(u.Host, h).ListenAndServe()
}

func makeServer(addr string, h http.Handler) *graceful.Server {
	if _, p, _ := net.SplitHostPort(addr); p != "" {
		log.Printf("server: must be on port :%s", p)
	}

	return &graceful.Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: h,
		},
		Timeout: 5 * time.Second,
		Logger:  log.New(os.Stderr, "server: ", 0)}
}
