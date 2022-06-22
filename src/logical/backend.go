package logical

import (
	"context"
	"github.com/36625090/involution/authorities"
	"github.com/go-various/consul"
	"github.com/go-various/redisplus"
	"github.com/go-various/xorm"
	"github.com/hashicorp/go-hclog"
)

// BackendContext 主配置
// 此处为了省事将 XormConfig RedisConfig  AuthSettings的设置放在了一个对象
type BackendContext struct {
	Logger       hclog.Logger
	Application  string					  `json:"application" hcl:"application"`
	XormConfig   *xorm.Config             `json:"xorm" hcl:"xorm,block"`
	RedisConfig  *redisplus.Config        `json:"redis" hcl:"redis,block"`
	AuthSettings *authorities.Settings    `json:"authorization" hcl:"authorization,block"`
	Consul       consul.Client            `json:"-"`
	TokenHandler authorities.TokenHandler `json:"-"`
}

func (m *BackendContext) Clone() *BackendContext {
	return &BackendContext{
		Application:  m.Application,
		Logger:       m.Logger,
		Consul:       m.Consul,
		XormConfig:   m.XormConfig,
		RedisConfig:  m.RedisConfig,
		AuthSettings: m.AuthSettings,
	}
}

type Factory func(context.Context, string, *BackendContext) (Backend, error)

//Backend 后端逻辑主入口
type Backend interface {
	//Initialize 初始化方法
	Initialize(context.Context) error
	//HandleRequest 处理客户端请求
	HandleRequest(context.Context, *Args) (*Reply, *WrapperError)
	//Documents 返回后端逻辑端点的文档结构
	Documents(context.Context) (*DocumentsReply, error)
	//Cleanup 清理函数
	Cleanup(context.Context)
	//BackendName 返回后端服务的名称
	BackendName() string

	//BackendDescription 返回后端服务的名称
	BackendDescription() string
}
