package client

import (
	"context"
	"fmt"
	"github.com/lsm1998/tinygo"
	mRpc "github.com/lsm1998/tinygo/pkg/server/middleware/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

const DailTimeout = 3 * time.Second

type connectMethod int

const (
	cmRegister connectMethod = 1
	cmDefault  connectMethod = 2
	cmProxy    connectMethod = 3
)

type (
	Client struct {
		conn *grpc.ClientConn
		Cli  interface{}
		opt  Options
		conf RpcClientConf
	}

	Option func(options *Options)

	Options struct {
		Timeout       time.Duration
		DialOptions   []grpc.DialOption
		Builder       resolver.Builder
		Debug         bool
		ProxyAddr     string
		connectMethod connectMethod
	}
)

func WithDebug() Option {
	return func(options *Options) {
		options.Debug = true
	}
}

func NewClient(c RpcClientConf, f func(*grpc.ClientConn) interface{}, opts ...Option) *Client {
	client := newClient(opts...)
	client.conf = c
	client.updateConnectMethod()
	client.updateDialOptions()
	client.createConn()
	client.Cli = f(client.conn)
	return client
}

// TODO 提供一个接口实现 dial
func (c *Client) dialByDns(addr string, serviceName string) (*grpc.ClientConn, error) {
	conn, err := c.dial(addr)
	if err != nil {
		return nil, fmt.Errorf("rpc dialByDns: %s, error: %s, make sure rpc service %s is already started",
			addr, err.Error(), serviceName)
	}
	return conn, nil
}

func (c *Client) dialByProxy() (*grpc.ClientConn, error) {
	return c.dial(c.opt.ProxyAddr)
}

func (c *Client) dialByBuilder(builder resolver.Builder, serviceName string) (*grpc.ClientConn, error) {
	target := fmt.Sprintf("%s://%s/%s/%s", builder.Scheme(), "docer_discov", tinygo.AppMod, serviceName)
	if c.opt.Debug {
		fmt.Println("target", target)
	}
	conn, err := c.dial(target)
	if err != nil {
		return nil, fmt.Errorf("rpc dialByBuilder: %s, error: %s, make sure rpc service %s is already started",
			target, err.Error(), serviceName)
	}
	return conn, nil
}

func (c *Client) dial(target string) (*grpc.ClientConn, error) {
	timeCtx, cancel := context.WithTimeout(context.Background(), DailTimeout)
	defer cancel()
	conn, err := grpc.DialContext(timeCtx, target, c.opt.DialOptions...)
	return conn, err
}

func WithBuilder(builder resolver.Builder) Option {
	return func(options *Options) {
		options.Builder = builder
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.Timeout = timeout
	}
}

// WithProxy 创建客户端的时候在metadata指定需要代理的grpc服务的应用域名和监听端口，即可在本地连接到kae上的grpc服务
// 调用该函数将默认使用 "kae-rpc-agent.wps.cn:81"
// from: liaoyijun
func WithProxy(addrs ...string) Option {
	return func(options *Options) {
		if tinygo.AppMod == "prod" {
			return
		}
		if len(addrs) > 0 {
			options.ProxyAddr = addrs[len(addrs)-1]
		} else {
			options.ProxyAddr = "kae-rpc-agent.wps.cn:81"
		}
	}
}

func WithMiddleware(interceptor grpc.UnaryClientInterceptor) Option {
	return func(options *Options) {
		options.DialOptions = append(options.DialOptions, grpc.WithChainUnaryInterceptor(interceptor))
	}
}

func (c *Client) updateConnectMethod() {
	if c.opt.ProxyAddr != "" {
		c.opt.connectMethod = cmProxy
	} else if c.opt.Builder != nil {
		c.opt.connectMethod = cmRegister
	} else {
		c.opt.connectMethod = cmDefault
	}
}

func (c *Client) createConn() {
	var err error
	switch c.opt.connectMethod {
	case cmRegister:
		c.conn, err = c.dialByBuilder(c.opt.Builder, c.conf.ServiceName)
	case cmDefault:
		c.conn, err = c.dialByDns(c.conf.Addr, c.conf.ServiceName)
	case cmProxy:
		c.conn, err = c.dialByProxy()
	}
	if err != nil {
		log.Fatal(err)
	}
}
func (c *Client) updateDialOptions() {
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(
			mRpc.TimeoutMiddleware(c.opt.Timeout),
		),
	}
	c.opt.DialOptions = append(options, c.opt.DialOptions...)
	switch c.opt.connectMethod {
	case cmRegister:
		s := fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, roundrobin.Name)
		c.opt.DialOptions = append([]grpc.DialOption{grpc.WithDefaultServiceConfig(s)}, c.opt.DialOptions...)
	case cmProxy:
		c.opt.DialOptions = append([]grpc.DialOption{grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(metadata.NewOutgoingContext(ctx, metadata.Pairs("target-endpoint", c.conf.Addr)), method, req, reply, cc, opts...)
		})}, c.opt.DialOptions...)
	}
}

func newClient(opts ...Option) *Client {
	var c = &Client{}
	for _, opt := range opts {
		opt(&c.opt)
	}
	return c
}
