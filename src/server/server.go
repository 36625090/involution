package server

import (
	"context"
	"fmt"
	"github.com/36625090/involution/authorities"
	"github.com/36625090/involution/config"
	"github.com/36625090/involution/logical"
	"github.com/36625090/involution/option"
	"github.com/36625090/involution/transport"
	"github.com/gin-gonic/gin"
	"github.com/go-various/consul"
	"github.com/hashicorp/go-hclog"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Server struct {
	sync.Mutex

	ctx           context.Context
	logger        hclog.InterceptLogger
	opts          *option.Options
	authorization authorities.Authorization
	httpTransport *transport.Transport
	httpServer    *http.Server
	netListener   net.Listener
	connection    *Connection
	backends      map[string]logical.Backend
	service       *consul.Service
	consulClient  consul.Client
	globalConfig  *config.GlobalConfig

}

func (m *Server) AddLoggerSinks(sinks ...hclog.SinkAdapter) {
	for _, sink := range sinks {
		m.logger.RegisterSink(sink)
	}
}

func NewServer(opts *option.Options, cfg *config.GlobalConfig ,cl consul.Client, logger hclog.InterceptLogger) *Server {

	return &Server{
		ctx:              context.Background(),
		globalConfig:	  cfg,
		opts:             opts,
		consulClient:     cl,
		logger:           logger,
		connection:       &Connection{},
		backends:         map[string]logical.Backend{},
	}
}

func (m *Server) Initialize() error {
	gin.SetMode(gin.ReleaseMode)
	en := gin.New()
	en.Use(gin.Recovery())

	m.httpTransport = transport.NewTransport(en, m.globalConfig.Transport, m.logger)
	if err := m.initContext(); err != nil {
		return err
	}

	m.initBackendAPIServer()

	if m.opts.Ui {
		m.addDocumentSchema()
		m.addDocumentUI()
	}

	return nil
}

//RegisterAuthorization 注册验证代理接口，如不需要课不注册
func (m *Server) RegisterAuthorization(authorization authorities.Authorization) error {
	m.authorization = authorization
	return nil
}

//Start the Server
//启动服务
func (m *Server) Start() error {
	m.logger.Info("start starting")
	addr := fmt.Sprintf("%s:%d", m.opts.Http.Address, m.opts.Http.Port)
	m.logger.Info("server listening on ", "address", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	m.netListener = l

	m.httpServer = &http.Server{
		Addr:         addr,
		Handler:      m.httpTransport,
		IdleTimeout:  time.Second * time.Duration(m.opts.Http.IdleTimeout),
		ReadTimeout:  time.Second * time.Duration(m.opts.Http.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(m.opts.Http.WriteTimeout),
	}

	go func() {
		if err := m.httpServer.Serve(m.netListener); err != nil && err != http.ErrServerClosed {
			log.Fatal("start http server: ", err)
			return
		}
	}()

	if m.consulClient != nil{
		if err := m.registerService(m.opts.Profile); err != nil {
			return err
		}
	}

	m.logger.Info("server start completed")
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	m.logger.Info("Server shutting")

	if m.consulClient != nil{
		m.consulClient.DeRegister(m.service)
	}

	if err := m.httpServer.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown:", err)
		return err
	}
	m.logger.Info("server shutdown completed")
	m.Cleanup()

	return nil
}

func (m *Server) Cleanup() {
	m.httpServer.Close()
	m.netListener.Close()
	for _, backend := range m.backends {
		backend.Cleanup(context.Background())
	}
}
