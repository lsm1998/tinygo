package server

import (
	"context"
	"fmt"
	"github.com/lsm1998/tinygo/pkg/logx"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"syscall"
	"time"
)

type Option func(o *Options)

type Options struct {
	sigs      []os.Signal
	servers   []Server
	delayTime time.Duration
	conf      Config
}

type Config struct {
	ServiceName   string     `yaml:"service_name"`
	Debug         bool       `yaml:"debug"`
	Env           string     `yaml:"env"`
	LogLevel      logx.Level `yaml:"log_level"`
	DeadDelayTime int        `yaml:"dead_delay_time"`
}

func Run(opts ...Option) {
	options := Options{
		sigs:      []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		delayTime: 3,
	}
	for _, o := range opts {
		o(&options)
	}
	options.run()
}

func WithDelayTime(delay int) Option {
	return func(o *Options) { o.delayTime = time.Duration(delay) }
}

func WithServer(srv ...Server) Option {
	return func(o *Options) { o.servers = append(o.servers, srv...) }
}

func (o *Options) run() {
	defer func() {
		if e := recover(); e != nil {
			log.Panic(e)
		}
	}()
	ctx, cancelFunc := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	for _, srv := range o.servers {
		func(srv Server) {
			eg.Go(func() error {
				<-ctx.Done() // wait for stop signal
				_ctx, _cancelFunc := context.WithTimeout(context.Background(), o.delayTime*time.Second)
				defer _cancelFunc()
				return srv.Shutdown(_ctx)
			})
			go func() {
				defer func() {
					if e := recover(); e != nil {
						log.Panic(e)
					}
				}()
				srv.Start()
			}()
		}(srv)
	}
	sig := <-Signal()
	fmt.Println("the received signal is %s", sig.String())
	cancelFunc()
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(fmt.Sprintf("service failed to run: %+v", err))
	}
}
