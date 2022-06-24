package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/36625090/involution/authorities"
	"github.com/36625090/involution/logical"
	"github.com/36625090/involution/logical/codes"
	"github.com/36625090/involution/transport"
	"github.com/36625090/involution/utils"
	"strings"
)

func (m *Server) RegisterBackend(bkName string, factory logical.Factory, cfg *logical.BackendContext) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.backends[bkName]; ok {
		return fmt.Errorf("existing backend: %s", bkName)
	}

	backend, err := factory(m.ctx, bkName, cfg)
	if err != nil {
		m.logger.Error("register backend", "name", bkName, "err", err)
		return err
	}

	if err := backend.Initialize(context.Background()); err != nil {
		m.logger.Error("initialize backend", "name", bkName, "err", err)
		return err
	}

	m.backends[bkName] = backend
	m.logger.Info("register backend", "name", bkName)
	return nil
}


func (m *Server) initBackendAPIServer() {

	path := strings.Join([]string{m.opts.Http.Path, "api"}, "/")
	if m.opts.Http.Trace {
		m.httpTransport.Use(m.loggerTracker(path))
	}

	m.httpTransport.AddHandle(path, logical.HttpMethodPOST, func(ctx *transport.Context) (err error) {
		request := ctx.Request()
		m.connection.Inc()
		defer func() {
			m.connection.Dec()
			if err != nil{
				m.connection.Error()
			}
		}()

		backend, ok := m.backends[request.Backend()]

		if !ok {
			ctx.WithCode(codes.CodeBackendIssue).WithMessage("invalid backend")
			return errors.New("invalid backend")
		}

		authorized, err := m.preAuthorization(request.Method, ctx.GetAuthToken())
		if err != nil {
			ctx.WithCode(codes.CodeUnauthorized).WithError(err)
			return err
		}

		args, err := ctx.DecodeArgs()
		if err != nil {
			ctx.WithCode(codes.CodeFailedDecodeArgs).WithError(err)
			return err
		}

		args.Authorized = authorized
		resp, werr := backend.HandleRequest(context.Background(), args)
		if werr != nil {
			ctx.WithCode(werr.Code).WithMessage(werr.String())
			return werr.Error()
		}
		if resp.Code != 0{
			ctx.WithCode(codes.ReturnCode(resp.Code)).WithMessage(resp.Message)
			return nil
		}
		resp.Pagination = nil
		ctx.WithContent(resp.Data)
		ctx.WithPagination(resp.Pagination)
		return nil
	})

}

func (m *Server) preAuthorization(method string, token string) (*authorities.Authorized, error) {
	if nil == m.authorization {
		return nil, errors.New("authorization unavailable")
	}

	if m.authorization.Settings().DefaultPolicy == authorities.AuthorizationPolicyAllow {
		return nil, nil
	}

	if utils.Contains(m.authorization.Settings().AnonMethods, method) {
		return nil, nil
	}

	if token == "" {
		return nil, errors.New("invalid token")
	}

	return m.authorization.Authentication(context.TODO(), token)
}
