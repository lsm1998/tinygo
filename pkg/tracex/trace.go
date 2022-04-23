package tracex

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"io"
	"net/http"
)

const KeySpan = "opentracing_span"

type Tracer struct {
	opentracing.Tracer
	io.Closer
	conf Config
}

type Span struct {
	opentracing.Span
	isDefault bool
}

func newSpanner(span opentracing.Span) Interface {
	return Span{Span: span}
}

//ValueSpan 从contex.Context中获取 span
//grpc控制器需要阅读 OpentracingGrpcServer 说明
//gin控制器需要阅读 OpentracingGrpcServer 说明
func ValueSpan(ctx context.Context) Interface {
	switch ctx.(type) {
	case *gin.Context:
		ctx = ctx.(*gin.Context).Request.Context()
	}
	value := ctx.Value(KeySpan)
	if value == nil {
		return newDefaultSpan()
	}
	span := value.(Span)
	return span
}

func (s Span) ChildContext(ctx context.Context) context.Context {
	switch ctx.(type) {
	case *gin.Context:
		ginCtx := ctx.(*gin.Context)
		ginCtx.Request = ginCtx.Request.WithContext(context.WithValue(ginCtx.Request.Context(), KeySpan, s))
		return ginCtx
	default:
		ctx = context.WithValue(ctx, KeySpan, s)
		return ctx
	}
}

func (s Span) IsRegister() bool {
	return s.isDefault
}

//InjectHttp 将span以header的形式注入到http.Request里
//这样子在http的接收方就可以解析头部的信息 ExtractHttp
func InjectHttp(req *http.Request, span Interface) error {
	return opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
}

func (s Span) ChildContextWithSpan(ctx context.Context, operationName string) context.Context {
	span := s.ChildSpan(operationName)
	return span.ChildContext(ctx)
}

//ChildSpan 在该Span延伸出子Span，形成一个完整的链路树
func (s Span) ChildSpan(operationName string) Interface {
	return Span{Span: opentracing.StartSpan(operationName, opentracing.ChildOf(s.Context()))}
}

//func (s Span) ChildContext(ctx context.Context, operationName string) context.Context {
//	return contextWithSpan(ctx, s.ChildSpan(operationName).(Span))
//}

// Info 跟s.LogKV是一样的，这里只是规范了下名字
func (s Span) Info(alternatingKeyValues ...interface{}) {
	alternatingKeyValues = append([]interface{}{"event", "info", "language", "golang"}, alternatingKeyValues...)
	s.LogKV(alternatingKeyValues...)
}

//Error 打印错误日志
func (s Span) Error(err error) {
	tag := s.SetTag("error", true)
	tag.LogFields(
		log.String("event", "error"),
		log.String("language", "golang"),
		log.String("err", err.Error()),
		log.String("stack", fmt.Sprintf("%+v", err)),
	)
}

func (s Span) TraceID() string {
	if !s.IsRegister() {
		return ""
	}
	return s.Context().(jaeger.SpanContext).TraceID().String()
}

func (s Span) Context() opentracing.SpanContext {
	return s.Span.Context()
}

func (s Span) Finish() {
	s.Span.Finish()
}

//ExtractHttp 解析HTTP请求头，导出Span
//这里会执行StartSpan命令，所以需要进行Finis操作
func ExtractHttp(req *http.Request) (Interface, error) {
	var s Span
	tracer := opentracing.GlobalTracer()
	extract, err := tracer.
		Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)
	if err == opentracing.ErrSpanContextNotFound {
		s.Span = tracer.StartSpan(req.URL.Path)
		return s, nil
	} else if err != nil {
		return s, err
	}
	s.Span = tracer.StartSpan(req.URL.Path, opentracing.ChildOf(extract))
	return s, nil
}
