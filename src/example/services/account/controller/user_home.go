package controller

import (
	"context"
	"github.com/36625090/involution/logical"
)

func (b *backend) userHome(ctx context.Context, args *logical.Args, reply *logical.Reply)*logical.WrapperError {
	reply.Data = map[string]interface{}{
		"name": "112230192c58",
	}
	return  nil
}

