package mRpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"runtime"
)

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 2048)
			buf = buf[:runtime.Stack(buf, true)]
			fmt.Printf("panic fail with error=[%v] stack==%s", e, buf)
		}
	}()

	return handler(ctx, req)
}
