package transport

import (
	"fmt"
	"github.com/36625090/involution/logical"
	"github.com/36625090/involution/logical/codes"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
	"runtime/debug"
)


//Transport 继承了gin实现的服务接口
type Transport struct {
	*gin.Engine
	logger hclog.Logger
	signer Signer
}

type Handle func(c *Context) error

func NewTransport(en *gin.Engine, settings *Settings ,logger hclog.Logger) *Transport {
	return &Transport{
		Engine: en,
		logger: logger.Named("transport"),
		signer: NewMD5Signer(settings),
	}
}

//AddHandle 添加路径handlerFunc
//path 绝对路径
func (m *Transport) AddHandle(absolutePath string, method logical.HttpMethod, handle Handle) {
	m.Engine.Handle(string(method), absolutePath, func(gCtx *gin.Context) {
		ctx := NewContext(gCtx)

		defer func() {
			if r := recover(); r != nil {
				m.logger.Error("received panic","err", r,"stack", string(debug.Stack()))
				ctx.WithCode(codes.CodeInvalidRequestParameter).
					WithMessage(fmt.Sprintf("%v", r))
				m.Write(ctx)
			}
		}()

		if err := ctx.ShouldBindJSON(); err != nil {
			m.logger.Error("should not bind JSON", "path", ctx.RawRequest().RequestURI, "err", err)

			ctx.WithCode(codes.CodeBindRequestData).
				WithMessage(err.Error())

			m.Write(ctx)
			return
		}

		if err := m.signer.Verify(ctx.GetClientID(), ctx.request.Sign, ctx.request); err != nil {
			m.logger.Error("verify request sign error", "path", ctx.RawRequest().RequestURI, "err", err)
			ctx.WithCode(codes.CodeInvalidSignature).
				WithMessage("verify request sign error, " + err.Error())

			m.Write(ctx)
			return
		}

		err := handle(ctx)

		m.logger.Error("handle request sign error",
			"path", ctx.RawRequest().RequestURI,"err", err)

		m.Write(ctx)

	})
}

func (m *Transport) Write(ctx *Context)  {
	sign, err := m.signer.Sign(GlobalSignKey, ctx.response)
	if err != nil {
		ctx.WithCode(codes.CodeInvalidSignature).WithMessage(err.Error())
		ctx.write()
		return
	}
	ctx.withSign(sign)
	ctx.write()
}

func (m *Transport) Router() gin.IRouter {
	return m.Engine
}
