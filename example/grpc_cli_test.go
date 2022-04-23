package example

import (
	"context"
	"fmt"
	"github.com/lsm1998/tinygo/pkg/etcdx"
	"github.com/lsm1998/tinygo/pkg/rpcx/client"
	"github.com/lsm1998/tinygo/pkg/rpcx/discov"
	"github.com/lsm1998/tinygo/pkg/tracex"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestGrpcCli(t *testing.T) {
	// no1 need registration center
	etcdClient := etcdx.Must(etcdx.WithConfig(etcdx.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}))

	// no2 init client
	RpcApiClient(client.RpcClientConf{ServiceName: "grpc_ser_test"}, etcdClient)

	// no3 ...
	for {
		time.Sleep(time.Second)
		resp, err := DefaultApiClient.SayHello(context.Background(), &SayHelloReq{Name: "lsm"})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(resp.Result)
	}

}

// apiClient maybe you have more than one client
type apiClient struct {
	DemoApiClient
}

var DefaultApiClient apiClient

func RpcApiClient(conf client.RpcClientConf, etcdClient *clientv3.Client) {
	client.NewClient(conf, func(conn *grpc.ClientConn) interface{} {
		DefaultApiClient.DemoApiClient = NewDemoApiClient(conn)
		return nil
	},
		client.WithMiddleware(tracex.OpentracingGrpcClient()),
		client.WithBuilder(discov.NewBuilder(etcdClient, conf.Addr)),
	)
	return
}
