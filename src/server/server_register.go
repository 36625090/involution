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

func (m *Server) unRegisterService() {
	if m.consulClient != nil {
		m.consulClient.DeRegister(m.service)
	}
}

//RegisterService 注册微服务
func (m *Server) registerService(tags ...string) error {
	if m.consulClient == nil {
		return nil
	}
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
		ID:             fmt.Sprintf("%s-%d-%d", m.opts.App, m.opts.Http.Port, rand.Int31()),
		Schema:         "http",
		Name:           m.opts.App,
		Address:        addr,
		MatchBody:      "",
		CheckInterval:  "30s",
		Port:           m.opts.Http.Port,
		Tags:           tags,
		HealthEndpoint: filepath.Join(m.opts.Http.Path, "/health"),
		ServiceAddress: map[string]api.ServiceAddress{
			consul.WanAddrKey: {Address: addr, Port: m.opts.Http.Port},
		},
	}
	return m.consulClient.Register(m.service)
}

func (m *Server) listenHealthyEndpoint() {
	path := filepath.Join(m.opts.Http.Path, "/health")
	m.logger.Trace("register health backend", "path", path)

	m.httpTransport.Handle("GET", path, func(c *gin.Context) {

		c.JSON(200, gin.H{
			"status":      "UP",
			"connections": m.connection,
			"memory":      utils.MemStats(),
			"cpus":        runtime.NumCPU(),
		})
	})
}
