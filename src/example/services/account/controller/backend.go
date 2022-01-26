package controller

import (
	"context"
	"github.com/36625090/involution/framework"
	"github.com/36625090/involution/logical"
)

func Factory(ctx context.Context, name string,conf *logical.BackendContext) (logical.Backend, error) {
	b := &backend{
		&framework.Backend{
			Config: conf,
			Name: name,
			Description: "账户管理",
		},
	}
	
	b.Endpoints =  framework.EndpointAppend(
		[]*framework.Endpoint{
			b.userPaths(),
		},
	)

	b.InitializeFunc = func(ctx context.Context) error {
		if b.LBAdapter != nil{
			b.LBAdapter.AddHooks(&microHook{Logger: b.Logger})
		}
		return nil
	}

	b.HandleRequestBeforeFunc = func(ctx context.Context, req *logical.Args) {

	}

	b.Clean = func(ctx context.Context) {

	}


	return b, nil
}

type backend struct {
	*framework.Backend
}
