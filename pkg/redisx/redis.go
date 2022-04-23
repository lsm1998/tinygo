package redisx

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Options struct {
	conf Config
}

type Option func(*Options)

type Config struct {
	Addr string `json:"addr" yaml:"addr"`
	Port int    `json:"port" yaml:"port"`
	Auth string `json:"auth" yaml:"auth"`
	Db   int    `json:"db" yaml:"db"`
}

func WithConfig(conf Config) Option {
	return func(k *Options) {
		k.conf = conf
	}
}

func Must(opts ...Option) *redis.Client {
	e := &Options{}
	for _, opt := range opts {
		opt(e)
	}
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", e.conf.Addr, e.conf.Port),
		Password: e.conf.Auth,
		DB:       e.conf.Db,
	})
}
