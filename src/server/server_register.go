package server

import (
	"fmt"
	"github.com/36625090/involution/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-various/consul"
	"github.com/hashicorp/consul/api"
	"math/rand"
	"path/filepath"
	"runtime"
)

//RegisterService 注册微服务
func (m *Server) registerService(tags ...string) error {
	addr := m.opts.Http.Address
	if addr == "" || addr == "0.0.0.0" {
		ip, err := utils.GetIP()
		if err != nil {
			m.logger.Error("get server address", "err", err)
			return err
		}
		addr = ip
	}

	m.service = &consul.Service{
		ID:             fmt.Sprintf("%s-%d-%d", m.opts.App, m.opts.Http.Port,rand.Int31()),
		Schema:         "http",
		Name:           m.opts.App,
		Address:        addr,
		MatchBody:      "",
		CheckInterval:  "30s",
		Port:           m.opts.Http.Port,
		Tags:           tags,
		HealthEndpoint: filepath.Join(m.opts.Http.Path, "health"),
		ServiceAddress: map[string]api.ServiceAddress{
			consul.WanAddrKey: {Address: addr, Port: m.opts.Http.Port},
		},
	}
	m.logger.Trace("register backend", "name", utils.JSONPrettyDump(m.service))

	m.httpTransport.Handle("GET", filepath.Join(m.opts.Http.Path, "health") , func(c *gin.Context) {

		c.JSON(200, gin.H{
			"status":      "UP",
			"connections": m.connection,
			"memory":      utils.MemStats(),
			"cpus":        runtime.NumCPU(),
		})
	})

	return m.consulClient.Register(m.service)
}
