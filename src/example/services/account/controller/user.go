package controller

import (
	"github.com/36625090/involution/example/services/account/views"
	"github.com/36625090/involution/framework"
	"github.com/36625090/involution/logical"
	"reflect"
)

const (
	auth		= "auth"
	home        = "home"
	login       = "login"
	logout      = "logout"
)

func (b *backend) userPaths() *framework.Endpoint {
	return &framework.Endpoint{
		Pattern:     "user",
		Description: "用户管理协议",
		Operations: map[string]framework.OperationHandler{
			login: &framework.EndpointOperation{
				Description: "用户登录",
				Callback:    b.userLogin,
				Input:       reflect.TypeOf(views.User{}),
				Output:      reflect.TypeOf(views.LoginReply{}),
			},

			home: &framework.EndpointOperation{
				Description: "用户主页",
				Callback:    b.userHome,
				Input:       reflect.TypeOf(logical.EmptyDocuments{}),
				Output:      reflect.TypeOf(logical.EmptyDocuments{}),
			},

			logout: &framework.EndpointOperation{
				Description: "用户登出",
				Callback:    b.userLogout,
				Input:       reflect.TypeOf(logical.EmptyDocuments{}),
				Output:      reflect.TypeOf(logical.EmptyDocuments{}),
			},

		},
	}
}
