package involution

import (
	"errors"
	"github.com/36625090/involution/authorities"
	"github.com/36625090/involution/config"
	"github.com/36625090/involution/logging"
	"github.com/36625090/involution/logical"
	"github.com/36625090/involution/option"
	"github.com/36625090/involution/server"
	"github.com/36625090/involution/utils"
	"github.com/go-various/consul"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
)



type Involution interface {
	//Initialize 服务初始化
	Initialize() error

	//RegisterBackend 注册后端逻辑端点
	RegisterBackend(string, logical.Factory, *logical.BackendContext)error

	//RegisterAuthorization 注册验证接口，如微注册则不验证
	RegisterAuthorization(authorization authorities.Authorization)error

	//Start 启动服务
	Start() error

	//AddLoggerSinks register sinks
	AddLoggerSinks(sinks ...hclog.SinkAdapter)
}

//DefaultInvolution default
func DefaultInvolution(opts *option.Options, factories  map[string]logical.Factory) (Involution, error) {
	if opts.Pprof{
		go func() {
			log.Println(http.ListenAndServe(opts.PprofAddr, nil))
		}()
	}

	logger, err :=  logging.NewLogger(opts.App, opts.Log)
	if err != nil {
		return nil, err
	}

	globalConfig := &config.GlobalConfig{}

	consulClient, err := initializeConfig(opts, globalConfig, logger)
	if err != nil {
		return nil, err
	}

	logger.Trace("initialize config", "config", utils.JSONPrettyDump(globalConfig))

	authorization, err := initializeAuthorization(globalConfig, err)
	if err != nil {
		return nil, err
	}

	context := &logical.BackendContext{
		Logger:       logger,
		Application:  opts.App,
		XormConfig:   globalConfig.XormConfig,
		RedisConfig:  globalConfig.RedisConfig,
		AuthSettings: globalConfig.Authorization,
		Consul:       consulClient,
		TokenHandler: authorization.TokenHandler(),
	}

	inv := server.NewServer(opts, globalConfig, consulClient, logger)
	if err := inv.RegisterAuthorization(authorization); err != nil {
		return nil, err
	}

	for name, factory := range factories {
		if err := inv.RegisterBackend(name, factory, context); err != nil {
			return nil, err
		}
	}

	if err := inv.Initialize(); err != nil {
		return nil, err
	}

	return inv, nil
}

func initializeConfig(opts *option.Options, globalConfig *config.GlobalConfig, logger hclog.Logger) (consul.Client, error) {
	if opts.ConfigFile != "" {
		bs, err := ioutil.ReadFile(opts.ConfigFile)
		if err != nil {
			return nil, err
		}
		if err := hcl.Unmarshal(bs, globalConfig); err != nil {
			return nil, err
		}
		return nil, nil
	}

	consulClient, err := consul.NewClient(opts.ConsulConfig(), logger)
	if err != nil {
		return nil, err
	}
	if err := consulClient.LoadConfig(globalConfig); err != nil {
		return nil, err
	}
	return consulClient, nil
}

func initializeAuthorization(globalConfig *config.GlobalConfig, err error) (authorities.Authorization, error) {
	var tokenHandler authorities.TokenHandler
	if globalConfig.Authorization.AuthType == authorities.AuthTypeJwt || globalConfig.Authorization.AuthType == "" {
		tokenHandler, err = authorities.NewJwtTokenHandler(globalConfig.Authorization)
	} else {
		tokenHandler, err = authorities.NewRedisTokenHandler(globalConfig.Authorization, globalConfig.RedisConfig)
	}
	if err != nil {
		return nil, err
	}

	authorization, err := authorities.NewAuthorization(globalConfig.Authorization, tokenHandler)
	if err != nil {
		return nil,  errors.New("initialization authorization: " + err.Error())
	}
	return  authorization, nil
}
