package grpcServer

import (
	"context"
	"fmt"
	"github.com/lsm1998/tinygo/pkg/rpcx"
	mRpc "github.com/lsm1998/tinygo/pkg/server/middleware/rpc"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

const (
	// DialTimeout the timeout of create connection
	DialTimeout = 5 * time.Second

	// BackoffMaxDelay provided maximum delay when backing off after failed connection attempts.
	BackoffMaxDelay = 3 * time.Second

	// KeepAliveTime is the duration of time after which if the client doesn't see
	// any activity it pings the server to see if the transport is still alive.
	KeepAliveTime = time.Duration(10) * time.Second

	// KeepAliveTimeout is the duration of time for which the client waits after having
	// pinged for keepalive check and if no activity is seen even after that the connection
	// is closed.
	KeepAliveTimeout = time.Duration(3) * time.Second

	// InitialWindowSize we set it 1GB is to provide system's throughput.
	InitialWindowSize = 1 << 30

	// InitialConnWindowSize we set it 1GB is to provide system's throughput.
	InitialConnWindowSize = 1 << 30

	// MaxSendMsgSize set max gRPC request message size sent to server.
	// If any request message size is larger than current value, an error will be reported from gRPC.
	MaxSendMsgSize = 4 << 30

	// MaxRecvMsgSize set max gRPC receive message size received from server.
	// If any message size is larger than current value, an error will be reported from gRPC.
	MaxRecvMsgSize = 4 << 30
)

type Server struct {
	server   *grpc.Server
	conf     Config
	handlers []grpc.UnaryServerInterceptor
	Register rpcx.Register
}

func newServer() *Server {
	return &Server{Register: rpcx.DefaultRegister}
}

func NewServer(f func(*Server), opts ...Option) *Server {
	s := newServer()
	for _, opt := range opts {
		opt(s)
	}
	//默认连接超时时间10s
	if s.conf.ConnectionTimeout <= 0 {
		s.conf.ConnectionTimeout = 10
	}
	timeoutOpt := grpc.ConnectionTimeout(time.Duration(s.conf.ConnectionTimeout) * time.Second)
	keepaliveParamsOpt := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(s.conf.KeepAliveConf.MaxConnectionIdle),
		MaxConnectionAge:      time.Duration(s.conf.KeepAliveConf.MaxConnectionAge),
		MaxConnectionAgeGrace: time.Duration(s.conf.KeepAliveConf.MaxConnectionAgeGrace),
		Time:                  time.Duration(s.conf.KeepAliveConf.Time),
		Timeout:               time.Duration(s.conf.KeepAliveConf.Timeout),
	})

	s.Use(mRpc.RecoveryInterceptor)
	var ops = []grpc.ServerOption{
		timeoutOpt, keepaliveParamsOpt,
		grpc.UnaryInterceptor(s.interceptor),
		grpc.InitialConnWindowSize(InitialConnWindowSize),
		grpc.InitialWindowSize(InitialWindowSize),
		grpc.MaxSendMsgSize(MaxSendMsgSize),
		grpc.MaxRecvMsgSize(MaxRecvMsgSize),
	}
	s.server = grpc.NewServer(ops...)
	f(s)
	for _, register := range registers[defaultRouterName] {
		register(s.Server())
	}
	return s
}

func (s *Server) Use(handlers ...grpc.UnaryServerInterceptor) *Server {
	s.handlers = append(s.handlers, handlers...)
	return s
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *Server) Start() {
	e := s.Register.Discov(s.conf.ServiceName, s.conf.Addr)
	go func() {
		for {
			select {
			case err := <-e:
				log.Errorf("grpc Server register %s failure %s\n", s.conf.Addr, err.Error())
			}
		}
	}()
	listen, err := net.Listen("tcp", s.conf.Addr)
	if err != nil {
		log.Fatal(fmt.Sprintf("start grpc Server failure %s", err.Error()))
		return
	}
	reflection.Register(s.server)
	logrus.Infof("start grpc Server listening %s", s.conf.Addr)
	if err = s.server.Serve(listen); err != nil {
		log.Fatal(fmt.Sprintf("grpc Server listening %s failure %s\n", s.conf.Addr, err.Error()))
	}
	logrus.Infof("grpc Server listening %s err: Server closed", s.conf.Addr)
	return
}

func (s *Server) interceptor(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var (
		i     int
		chain grpc.UnaryHandler
	)

	n := len(s.handlers)
	if n == 0 {
		return handler(ctx, req)
	}

	chain = func(ic context.Context, ir interface{}) (interface{}, error) {
		if i == n-1 {
			return handler(ic, ir)
		}
		i++
		return s.handlers[i](ic, ir, args, chain)
	}
	return s.handlers[0](ctx, req, args, chain)
}

func (s *Server) Server() *grpc.Server {
	return s.server
}
