package tracex

import (
	"github.com/opentracing/opentracing-go"
)

func newTracer() *Tracer {
	return &Tracer{}
}

func newDefaultSpan() Span {
	return Span{isDefault: true, Span: opentracing.GlobalTracer().StartSpan("default value span")}
}
