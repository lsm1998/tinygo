package tracex

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

// OpentracingGin 将在此中间件内生成的span赋值给ctx向下传递，
// 下层通过tracex.ValueSpan(ctx) 获取span无需再次进行Finish
func OpentracingGin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !opentracing.IsGlobalTracerRegistered() {
			ctx.Next()
			return
		}
		span, err := ExtractHttp(ctx.Request)
		defer span.Finish()
		if err != nil {
			logrus.Error(err)
			ctx.Abort()
			return
		}
		//将在此中间件内生成的span赋值给ctx向下传递，
		//下层通过tracex.ValueSpan(ctx) 获取span无需再次进行Finish
		span.ChildContext(ctx)
		ctx.Next()
	}
}
