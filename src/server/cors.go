package server

import (
	"github.com/36625090/involution/logical"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Cors() gin.HandlerFunc {
	headers := []string{
		"Origin", "Authorization", "Content-Type",
		string(logical.HeaderTraceIDKey), string(logical.HeaderApplicationKey),string(logical.HeaderClientIDKey),
		"Os-Version", "App-Version", "Location",
	}
	mwCORS := cors.New(cors.Config{
		//准许跨域请求网站,多个使用,分开,限制使用*
		AllowAllOrigins: true,
		//准许使用的请求方式
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		//准许使用的请求表头
		AllowHeaders: headers,
		//显示的请求表头
		ExposeHeaders: []string{"Content-Type"},
		//凭证共享,确定共享
		AllowCredentials: true,
		//超时时间设定
		MaxAge: 24 * time.Hour,
	})
	return mwCORS
}