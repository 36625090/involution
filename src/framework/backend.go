package framework

import (
	"context"
	"errors"
	"fmt"
	"github.com/36625090/involution/authorities"
	"github.com/36625090/involution/logical"
	"github.com/36625090/involution/logical/codes"
	"github.com/36625090/involution/utils"
	"github.com/go-playground/validator/v10"
	"github.com/go-various/consul"
	"github.com/go-various/micro"
	"github.com/go-various/redisplus"
	"github.com/go-various/xorm"
	log "github.com/hashicorp/go-hclog"
	"runtime/debug"
	"strings"
	"sync"
)

var ErrEndpointNotExists = errors.New("endpoint not found")
var ErrOperationNotExists = errors.New("operation not exists")

// Backend is an implementation of logical.Backend
var _ logical.Backend = (*Backend)(nil)

// Backend 后端主逻辑框架实现
type Backend struct {
	Name                    string
	Description             string
	Logger                  log.Logger
	Config                  *logical.BackendContext
	once                    sync.Once
	documents               logical.Documents
	Endpoints               []*Endpoint
	ConsulClient            consul.Client
	XormPlus                xorm.EngineInterface
	RedisCli                redisplus.RedisCli
	LBAdapter               micro.LBAdapter
	TokenHandler            authorities.TokenHandler
	Clean                   CleanupFunc
	InitializeFunc          InitializeFunc
	HandleRequestBeforeFunc HandleRequestBeforeFunc
	validator               *validator.Validate
}

func (b *Backend) BackendName() string {
	return b.Name
}
func (b *Backend) BackendDescription() string {
	return b.Description
}

// InitializeFunc 初始化函数类型
type InitializeFunc func(context.Context) error

//HandleRequestBeforeFunc 请求前的操作函数
type HandleRequestBeforeFunc func(context.Context, *logical.Args)

// OperationFunc 请求操作函数
type OperationFunc func(context.Context, *logical.Args, *logical.Reply) *logical.WrapperError

// CleanupFunc 清理函数
type CleanupFunc func(context.Context)

// HandleRequest 处理请求逻辑
func (b *Backend) HandleRequest(ctx context.Context, req *logical.Args) (resp *logical.Reply, err *logical.WrapperError) {
	if b.Logger.IsTrace() {
		b.Logger.Trace("handle request before", "request", utils.JSONDump(req))
	}
	resp = &logical.Reply{}
	defer func() {
		if b.Logger.IsTrace() {
			b.Logger.Trace("handle request after", "response", utils.JSONDump(resp))
		}

		if err != nil {
			b.Logger.Error("handle request",
				"backend", b.Config.Application, "endpoint",
				req.Endpoint, "operation", req.Operation,
				"error", err)
		}
		if err2 := recover(); nil != err2 {
			b.Logger.Error("recover panic", "err", err2, "stack", string(debug.Stack()))
			err = &logical.WrapperError{
				Code:  codes.CodeServerInternalError,
				Scope: "",
				Err:   err2.(error),
			}
		}
	}()

	if err := b.Validate(req); err != nil {
		return nil, &logical.WrapperError{
			Code:  codes.CodeDataValidateException,
			Scope: "",
			Err:   err,
		}
	}

	if req.GetTraceID() == "" {
		return nil, &logical.WrapperError{
			Code:  codes.CodeRequestHeaderMissing,
			Scope: "",
			Err:   errors.New("trace ID is required"),
		}
	}
	// Find the matching route
	path := b.find(req.Endpoint)
	if nil == path {
		return nil, &logical.WrapperError{
			Code:  codes.CodeEndpointNotFound,
			Scope: "",
			Err:   ErrEndpointNotExists,
		}
	}

	operation, ok := path.Operations[req.Operation]
	if !ok {
		return nil, &logical.WrapperError{
			Code:  codes.CodeOperationNotFound,
			Scope: "",
			Err:   ErrOperationNotExists,
		}
	}

	if operation.Handler() == nil {
		return nil, &logical.WrapperError{
			Code:  codes.CodeOperationHandlerIssue,
			Scope: "",
			Err: fmt.Errorf("operation headler: %s.%s.%s cannot be nil",
				req.Backend, req.Endpoint, req.Operation),
		}
	}

	if b.HandleRequestBeforeFunc != nil {
		b.HandleRequestBeforeFunc(ctx, req)
	}
	err = operation.Handler()(ctx, req, resp)
	return resp, err
}

// Cleanup 清理函数
func (b *Backend) Cleanup(ctx context.Context) {
	b.Logger.Trace("cleaning")
	if b.Clean != nil {
		b.Clean(ctx)
	}
}

// Initialize 框架初始化函数
func (b *Backend) Initialize(ctx context.Context) (err error) {
	b.validator = validator.New()

	b.Logger = b.Config.Logger.Named("backend").Named(b.Name)

	if err := b.checkEndpoint(); err != nil {
		return err
	}

	//初始化xorm
	b.XormPlus, err = xorm.NewEnginePlus(b.Config.XormConfig, b.Logger.StandardWriter(&log.StandardLoggerOptions{}))
	if err != nil {
		return err
	}

	//初始化redis
	prefix := strings.Join([]string{b.Config.Application, b.Name}, ":")
	b.RedisCli, err = redisplus.NewRedisCli(b.Config.RedisConfig, prefix)
	if err != nil {
		return err
	}

	//初始化微服务客户端
	if b.Config.Consul != nil {
		b.ConsulClient = b.Config.Consul
		b.LBAdapter = micro.RandomLBClient(logical.NewMicroServiceClient(b.ConsulClient))
	}

	//初始化验证接口
	b.TokenHandler = b.Config.TokenHandler

	if b.InitializeFunc != nil {
		b.InitializeFunc(ctx)
	}

	return err
}

//检查路径配置是否正确
func (b *Backend) checkEndpoint() error {

	for _, p := range b.Endpoints {
		for operation, handler := range p.Operations {
			if handler.Handler() == nil {
				return fmt.Errorf("operation callback: %s.%s.%s cannot be nil",
					b.Name, p.Pattern, operation)
			}
		}
	}
	return nil
}

//精准模式
func (b *Backend) find(path string) *Endpoint {
	for _, p := range b.Endpoints {
		if p.Pattern == path {
			return p
		}
	}
	return nil
}

// Documents 返回接口文档信息
func (b *Backend) Documents(ctx context.Context) (*logical.DocumentsReply, error) {
	b.once.Do(b.initDocumentsOnce)
	return &logical.DocumentsReply{
		Documents: b.documents,
	}, nil
}
