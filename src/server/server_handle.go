package server

import (
	"github.com/36625090/involution/config"
	"github.com/36625090/involution/logical"
	"github.com/36625090/involution/transport"
	"github.com/go-various/consul"
	"github.com/hashicorp/go-hclog"
)

type HandlerParams struct {
	Backends map[string]logical.Backend
	Config   *config.GlobalConfig
	Consul   consul.Client
	Logger   hclog.InterceptLogger
}

func (m *Server) AddHandle(path string, method logical.HttpMethod, handle func(*transport.Context, *HandlerParams) error) {
	m.httpTransport.AddHandle(path, method, func(c *transport.Context) error {
		params := &HandlerParams{
			Backends: m.backends,
			Config:   m.globalConfig,
			Consul:   m.consulClient,
			Logger:   m.logger,
		}
		return handle(c, params)
	})
}
