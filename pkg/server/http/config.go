package httpServer

import (
	"github.com/gin-gonic/gin"
	"github.com/lsm1998/tinygo"
	mHttp "github.com/lsm1998/tinygo/pkg/server/middleware/http"
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
	Name         string `yaml:"name"`
	Addr         string `yaml:"addr"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

func NewServer(f func(*Server), opts ...Option) *Server {
	server := newServer()
	for _, opt := range opts {
		opt(server)
	}
	if server.conf.Name == "" {
		server.conf.Name = defaultRouterName
	}
	server.server = &http.Server{
		Addr:           server.conf.Addr,
		ReadTimeout:    time.Duration(server.conf.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(server.conf.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.Engine = gin.New()
	if tinygo.AppMod == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		server.Use(gin.Logger())
	}
	server.Use(mHttp.Recovery())
	if server.withCors {
		server.Use(mHttp.Cross())
	}
	f(server)
	//TODO 这个目前仅支持一个默认的http路由配置,目前是map形式方便后续扩展
	for _, register := range registers[server.conf.Name] {
		register(server.Engine)
	}
	server.server.Handler = server.Engine
	return server
}

func WithoutCors() Option {
	return func(s *Server) {
		s.withCors = false
	}
}
