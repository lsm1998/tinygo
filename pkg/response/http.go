package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lsm1998/tinygo/pkg/response/ecode"
)

const (
	fail = "error"
	ok   = "ok"
)

type Response struct {
	Data    interface{} `json:"data"`
	Result  string      `json:"result"`
	TraceID string      `json:"trace_id,omitempty"`
	*ecode.ErrorX
}

func (r *Response) init() {
	r.Data = nil
	r.Result = ""
	r.ErrorX = nil
}

func (r *Response) withData(data interface{}) *Response {
	r.Data = data
	return r
}

func (r *Response) withResult(isSuccess bool) *Response {
	if isSuccess {
		r.Result = ok
	} else {
		r.Result = fail
	}
	return r
}

func (r *Response) withError(errno *ecode.ErrorX) *Response {
	r.ErrorX = errno
	if errno.Errno != ecode.Success && errno.Ext != nil {
		errno.ExtMessage = fmt.Sprintf("%v", errno.Ext)
	}
	return r
}

func Success(ctx *gin.Context, data interface{}) {
	getResponseInstance(ctx, successInterceptor(data))
}

//Error 返回错误信息到前端
func Error(ctx *gin.Context, err error, data ...interface{}) {
	getResponseInstance(ctx, errorInterceptor(err, data))
}
