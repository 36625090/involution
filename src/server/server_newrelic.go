package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
	"github.com/newrelic/go-agent/v3/newrelic"
	"io"
)

func (m *Server) newrelicTracer() gin.HandlerFunc {
	var l = io.Discard
	if m.opts.NewrelicTrace {
		l = m.logger.StandardWriter(&hclog.StandardLoggerOptions{})
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(m.opts.App),
		newrelic.ConfigLicense(m.opts.NewrelicKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigDebugLogger(l),
	)
	if nil != err {
		m.logger.Error("initial newrelic", "err", err)
		return nil
	}
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		txn := app.StartTransaction(c.Request.RequestURI)
		txn.SetWebRequestHTTP(c.Request)
		txn.SetWebResponse(c.Writer)
		c.Next()
		txn.End()
	}
}
