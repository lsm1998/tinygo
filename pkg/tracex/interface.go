package tracex

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

// Opentracer 定义公共接口，支持Zipkin Jaeger等等，目前仅实现 Jeager
type Opentracer interface {
	TraceID() string
}

type Interface interface {
	Info(alternatingKeyValues ...interface{})
	Error(err error)
	Context() opentracing.SpanContext
	Finish()
	IsRegister() bool
	ChildSpan(operationName string) Interface
	ChildContext(ctx context.Context) context.Context
	ChildContextWithSpan(ctx context.Context, operationName string) context.Context
	Opentracer
}
