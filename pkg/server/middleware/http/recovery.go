package mHttp

import (
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				buf := make([]byte, 2048)
				buf = buf[:runtime.Stack(buf, true)]
				log.Panicf("panic fail with error=[%v] stack==%s", e, buf)
			}
		}()
		c.Next()
	}
}
