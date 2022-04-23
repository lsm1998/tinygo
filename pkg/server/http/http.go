package httpServer

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	server *http.Server
	*gin.Engine
	conf     Config
	withCors bool
}

func newServer() *Server {
	return &Server{withCors: true}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Use(middleware ...gin.HandlerFunc) {
	s.Engine.Use(middleware...)
}

func (s *Server) Start() {
	logrus.Infof("start http server listening %s", s.conf.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		logrus.Errorf("http server listening %s err: %s \n", s.conf.Addr, err)
	}
}
