package server

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/tylerb/graceful"
)

var serverName = filepath.Base(os.Args[0])

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
		log.Printf("%s is now ready to accept connections on port %s", serverName, p)
	}

	return &graceful.Server{
		Server: &http.Server{
			Addr:    host,
			Handler: h,
		},
		Timeout: 5 * time.Second,
		Logger:  log.New(os.Stderr, serverName+" ", 0)}
}
