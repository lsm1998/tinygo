package websocket

import (
	httpServer "github.com/lsm1998/tinygo/pkg/server/http"
	"net/http"
	"time"
)

type Option func(c *Server)

func SetConf(c Config) Option {
	return func(s *Server) {
		s.conf = c
	}
}

type Config struct {
	Path              string `json:"path" yaml:"path"`
	httpServer.Config `json:",inline" yaml:",inline"`
}

func NewServer(f func(*Server), opts ...Option) *Server {
	server := newServer()
	for _, opt := range opts {
		opt(server)
	}
	server.server = &http.Server{
		Addr:           server.conf.Addr,
		ReadTimeout:    time.Duration(server.conf.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(server.conf.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	f(server)
	return server
}
