package tracex

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TextMD struct {
	metadata.MD
}

// Set 重写TextMapWriter的Set方法，我们需要将carrier中的数据写入到metadata中，这样grpc才会携带。
func (t TextMD) Set(key, val string) {
	t.MD[key] = append(t.MD[key], val)
}

// ForeachKey 读取metadata中的span信息
func (t TextMD) ForeachKey(handler func(key, val string) error) error { //不能是指针
	for key, val := range t.MD {
		for _, v := range val {
			if err := handler(key, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func fromIncomingContext(ctx context.Context) TextMD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}
	textMD := TextMD{MD: md}
	return textMD
}

func fromOutgoingContext(ctx context.Context) TextMD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}
	textMD := TextMD{MD: md}
	return textMD
}

// OpentracingGrpcClient 客户端解码span数据并组装成可以再GRPC中解析的数据
func OpentracingGrpcClient() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if !opentracing.IsGlobalTracerRegistered() {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		span := ValueSpan(ctx)
		tracer := opentracing.GlobalTracer()
		textMD := fromOutgoingContext(ctx)
		err := tracer.Inject(span.Context(), opentracing.TextMap, textMD)
		if err != nil {
			return err
		}
		ctx = metadata.NewOutgoingContext(ctx, textMD.MD)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// OpentracingGrpcServer 将在此中间件内生成的span赋值给ctx向下传递，
// 下层通过tracex.ValueSpan(ctx) 获取span无需再次进行Finish
func OpentracingGrpcServer() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if !opentracing.IsGlobalTracerRegistered() {
			return handler(ctx, req)
		}
		tracer := opentracing.GlobalTracer()
		textMD := fromIncomingContext(ctx)
		spanContext, err := tracer.Extract(opentracing.TextMap, textMD)
		if err != nil {
			if err == opentracing.ErrSpanContextNotFound {
				return handler(ctx, req)
			}
			return nil, err
		}
		span := tracer.StartSpan(info.FullMethod, opentracing.ChildOf(spanContext))
		span.SetTag("protocol", "grpc")
		span.SetTag("method", info.FullMethod)
		defer span.Finish()
		//将在此中间件内生成的span赋值给ctx向下传递，
		//下层通过tracex.ValueSpan(ctx) 获取span无需再次进行Finish
		spanner := newSpanner(span)
		ctx = spanner.ChildContext(ctx)
		return handler(ctx, req)
	}
}
