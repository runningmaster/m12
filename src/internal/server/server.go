package server

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tylerb/graceful"
)

// Run starts a server with router
func Run(addr string, h http.Handler) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	return makeServer(u.Host, h).ListenAndServe()
}

func makeServer(host string, h http.Handler) *graceful.Server {
	if _, p, _ := net.SplitHostPort(host); p != "" {
		log.Printf("server: is now ready to accept connections on port %s", p)
	}

	return &graceful.Server{
		Server: &http.Server{
			Addr:    host,
			Handler: h,
		},
		Timeout: 5 * time.Second,
		Logger:  log.New(os.Stderr, "server: ", 0)}
}
