package response

import (
	"github.com/gin-gonic/gin"
	"github.com/lsm1998/tinygo"
	"github.com/lsm1998/tinygo/pkg/response/ecode"
	"github.com/lsm1998/tinygo/pkg/tracex"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
	"sync"
)

type InterceptorFunc func(ctx *gin.Context, p *Response)

var (
	respPool = &sync.Pool{
		New: func() interface{} {
			return &Response{}
		},
	}
)

func getResponseInstance(ctx *gin.Context, interceptors ...InterceptorFunc) {
	p := respPool.Get().(*Response)
	for _, interceptor := range interceptors {
		interceptor(ctx, p)
	}
	traceInterceptor(ctx, p)
	returnValue(ctx, p)
	p.init()
	respPool.Put(p)
}

// 实现默认的返回值拦截器
var returnValue = InterceptorFunc(func(ctx *gin.Context, p *Response) {
	ctx.JSON(200, p)
})

// 实现默认的tracex拦截器，注入 trace_id
var traceInterceptor = InterceptorFunc(func(ctx *gin.Context, p *Response) {
	if tinygo.AppTrace == "true" {
		p.TraceID = tracex.ValueSpan(ctx).TraceID()
	}
})

// 实现接口返回success拦截器
func successInterceptor(data interface{}) InterceptorFunc {
	return InterceptorFunc(func(ctx *gin.Context, p *Response) {
		p = p.withData(data).withResult(true).withError(ecode.Wrap(ecode.Success, nil))
	})
}

// 实现接口返回error拦截器，注入错误值
func errorInterceptor(err error, data ...interface{}) InterceptorFunc {
	return InterceptorFunc(func(ctx *gin.Context, p *Response) {
		p = p.withResult(false)
		if len(data) == 1 {
			p = p.withData(data[0])
		} else {
			p = p.withData(data)
		}

		if e1, ok := err.(*ecode.ErrorX); ok {
			p = p.withError(e1)
		} else if e2, ok := err.(*ecode.Errno); ok {
			p = p.withError(ecode.Wrap(e2, nil))
		} else {
			originErr := errors.Cause(err)
			if rpcStatus, ok := status.FromError(originErr); ok {
				// grpc err错误
				var rpcErr ecode.Errno
				msgData := strings.Split(rpcStatus.Message(), "<sp>")
				if len(msgData) == 2 {
					msgCode, e := strconv.Atoi(msgData[1])
					if e == nil {
						rpcErr.Code = msgCode
					}
				}
				rpcErr.Msg = msgData[0]
				p = p.withError(ecode.Wrap(&rpcErr, nil))
			} else {
				p = p.withError(ecode.Wrap(ecode.Error, err))
			}
		}
	})
}
