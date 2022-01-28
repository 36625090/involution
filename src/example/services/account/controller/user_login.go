package controller

import (
	"context"
	"encoding/json"
	"github.com/36625090/involution/authorities"
	"github.com/36625090/involution/example/services/account/model"
	"github.com/36625090/involution/example/services/account/views"
	"github.com/36625090/involution/logical"
)

func (b *backend) userLogin(ctx context.Context, args *logical.Args, reply *logical.Reply) *logical.WrapperError  {

	authorized := &authorities.Authorized{
		ID:           "1",
		Account:      "example-account",
		AccountRoles:  []string{"member"},
		Principal: map[string]interface{}{
			"mobile": "13800000000",
		},
	}
	user := &views.User{}
	if err := args.ShouldBindJSON(user); err != nil {
		reply.Code = 100
		reply.Message = err.Error()
		return nil
	}

	dbUser := model.User{
		Mobile: new(string),
	}
	*dbUser.Mobile = user.Mobile
	has, err := b.XormPlus.Get(&dbUser)
	if err != nil {
		reply.Code = 101
		reply.Message = err.Error()
		return nil
	}
	if !has{
		reply.Code = 102
		reply.Message = "user not found"
		return nil
	}

	token, err := b.TokenHandler.GenerateToken(authorized)
	if err != nil {
		reply.Code = 103
		reply.Message = err.Error()
		return nil
	}

	cb, _ := json.Marshal(authorized)
	if err := b.RedisCli.Set(authorized.ID, cb, "1h"); err != nil {
		reply.Code = 104
		reply.Message = err.Error()
		return nil
	}
	resp := views.LoginReply{
		Token:     token,
		Principal: dbUser,
	}
	reply.Data = resp
	return nil
}