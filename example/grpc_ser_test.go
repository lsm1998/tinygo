package example

import (
	"context"
	"github.com/lsm1998/tinygo/pkg/etcdx"
	"github.com/lsm1998/tinygo/pkg/rpcx/discov"
	"github.com/lsm1998/tinygo/pkg/server"
	grpcServer "github.com/lsm1998/tinygo/pkg/server/grpc"
	"github.com/prometheus/common/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"testing"
)

func TestGrpcSer(t *testing.T) {
	// no1 need registration center
	etcdClient := etcdx.Must(etcdx.WithConfig(etcdx.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}))

	// no2 grpc api implementation
	grpcServer.AddComponent(func(server *grpc.Server) {
		RegisterDemoApiServer(server, ser)
	})

	// no3 server Run
	server.Run(
		server.WithServer(
			NewRPC(etcdClient),
		),
	)
	log.Info("server exiting...")
}

func NewRPC(etcdClient *clientv3.Client) server.Server {
	return grpcServer.NewServer(func(s *grpcServer.Server) {},
		grpcServer.SetConf(grpcServer.Config{
			Addr:        ":8081",
			ServiceName: "grpc_ser_test",
		}), grpcServer.SetRegister(discov.NewRegister(etcdClient)))
}

var ser = &exampleSer{}

type exampleSer struct{}

func (e exampleSer) SayHello(ctx context.Context, req *SayHelloReq) (*SayHelloRsp, error) {
	return &SayHelloRsp{
		Result: "hello:" + req.Name,
	}, nil
}
