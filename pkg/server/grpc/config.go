package grpcServer

import (
	"github.com/lsm1998/tinygo/pkg/rpcx"
	"go.etcd.io/etcd/client/v3"
)

type Option func(c *Server)

func SetConf(c Config) Option {
	return func(s *Server) {
		s.conf = c
	}
}

func SetRegister(register rpcx.Register) Option {
	return func(c *Server) {
		c.Register = register
	}
}

type Config struct {
	Addr              string        `yaml:"addr"`
	ServiceName       string        `yaml:"service_name"`
	ConnectionTimeout int           `yaml:"connection_timeout"` //所有连接的超时设置
	KeepAliveConf     keepAliveConf `yaml:"keep_alive_conf"`
	client            *clientv3.Client
	Register          `json:"-" yaml:"-"`
}

type keepAliveConf struct {
	MaxConnectionIdle     int `yaml:"max_connection_idle"` //
	MaxConnectionAge      int `yaml:"max_connection_age"`
	MaxConnectionAgeGrace int `yaml:"max_connection_age_grace"`
	Time                  int `yaml:"time"`
	Timeout               int `yaml:"timeout"`
}
