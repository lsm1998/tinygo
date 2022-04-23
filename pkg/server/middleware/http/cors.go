package mHttp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cross() gin.HandlerFunc {
	return func(r *gin.Context) {
		origin := r.Request.Header.Get("Origin")
		r.Header("Access-Control-Allow-Origin", origin)
		r.Header("Access-Control-Allow-Credentials", "true")
		r.Header("Access-Control-Allow-Headers", "COOKIE,token,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,appid,accept-language")
		r.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,HEAD")
		if r.Request.Method == "OPTIONS" {
			r.AbortWithStatus(http.StatusNoContent)
		}
		r.Next()
	}
}
