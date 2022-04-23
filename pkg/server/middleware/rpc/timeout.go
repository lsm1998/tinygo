package mRpc

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

func TimeoutMiddleware(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		if deadline, ok := ctx.Deadline(); ok {
			leftTime := time.Until(deadline)
			if leftTime < timeout {
				timeout = leftTime
			}
		}
		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(timeout))
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
