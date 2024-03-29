package etcdx

import (
	"context"
	"errors"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type Options struct {
	conf Config
}

type Option func(*Options)

// Config etcd的配置字段
type Config struct {
	Endpoints []string `yaml:"endpoints"`
	Timeout   int64    `yaml:"timeout"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password" kms:"encode"`
}

func WithConfig(conf Config) Option {
	return func(k *Options) {
		k.conf = conf
	}
}

func Optional(opts ...Option) *clientv3.Client {
	client, err := newClient(opts...)
	if err != nil {
		fmt.Println(err.Error())
	}
	return client
}

func Must(opts ...Option) *clientv3.Client {
	client, err := newClient(opts...)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	return client
}

func newClient(opts ...Option) (*clientv3.Client, error) {
	e := &Options{}
	for _, opt := range opts {
		opt(e)
	}
	return clientv3.New(clientv3.Config{
		Endpoints:   e.conf.Endpoints,
		DialTimeout: time.Duration(e.conf.Timeout) * time.Second,
		Username:    e.conf.Username,
		Password:    e.conf.Password,
	})
}

func MustEtcdWithTimeout(ctx context.Context, opts ...Option) (*clientv3.Client, error) {
	var (
		client *clientv3.Client
		err    error
		c      = make(chan struct{}, 1)
	)

	go func() {
		client, err = newClient(opts...)
		c <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		<-c
		log.Println("Timeout")
		return nil, errors.New("init etcd client timeout")
	case <-c:
		return client, err
	}
}
