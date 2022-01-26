package server

import (
	"context"
	"fmt"
	"github.com/36625090/involution/logical"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

//initContext 初始化gin服务中间件
func (m *Server) initContext() error {
	m.httpTransport.Use(m.requestTracer(context.Background()))

	if m.opts.Newrelic{
		m.httpTransport.Use(m.newrelicTracer())
	}

	if m.opts.Http.Cors {
		m.httpTransport.Use(Cors())
	}

	return nil
}

func (m *Server) requestTracer(ctx context.Context) gin.HandlerFunc {

	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			return
		}
		// 设置 trace-id 变量
		c.Request.Header.Set(string(logical.HeaderApplicationKey), m.opts.App)
		c.Request.Header.Set(string(logical.HeaderTraceIDKey), uuid.New().String())

		c.Next()
	}
}

func (m *Server) loggerTracker(path string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		t := time.Now()
		defer func() {
			// 请求后
			latency := time.Since(t)
			if c.Request.RequestURI != path{
				return
			}
			msg := fmt.Sprintf("client=%v client-id=%s trace-id=%s application=%s uri=%v code=%d latency=%v",
				c.Request.RemoteAddr,
				c.Request.Header.Get(logical.HeaderClientIDKey.String()),
				c.Request.Header.Get(logical.HeaderTraceIDKey.String()),
				c.Request.Header.Get(logical.HeaderApplicationKey.String()),
				c.Request.RequestURI, c.Writer.Status(), latency)
			m.logger.Info(msg)
		}()
		defer c.Next()

	}
}

func Cors() gin.HandlerFunc {
	headers := []string{
		"Origin", "Authorization", "Content-Type",
		string(logical.HeaderTraceIDKey), string(logical.HeaderApplicationKey),string(logical.HeaderClientIDKey),
		"Os-Version", "App-Version", "Location",
	}
	mwCORS := cors.New(cors.Config{
		//准许跨域请求网站,多个使用,分开,限制使用*
		AllowOrigins: []string{"*"},
		//准许使用的请求方式
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		//准许使用的请求表头
		AllowHeaders: headers,
		//显示的请求表头
		ExposeHeaders: []string{"Content-Type"},
		//凭证共享,确定共享
		AllowCredentials: true,
		//容许跨域的原点网站,可以直接return true就万事大吉了
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		//超时时间设定
		MaxAge: 24 * time.Hour,
	})
	return mwCORS
}