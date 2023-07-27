package websocket

import (
	"context"
	"github.com/lsm1998/tinygo/pkg/logx"
	"net/http"
)

type Server struct {
	mux      *http.ServeMux
	server   *http.Server
	conf     Config
	withCors bool
}

func newServer() *Server {
	return &Server{withCors: true, mux: http.NewServeMux()}
}

func (s *Server) HandleFunc(handle func(w http.ResponseWriter, r *http.Request)) {
	s.mux.HandleFunc(s.conf.Path, handle)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Start() {
	s.server.Handler = s.mux
	logx.Infof("start websocket server listening %s", s.conf.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		logx.Errorf("websocket server listening %s err: %s \n", s.conf.Addr, err)
	}
}
