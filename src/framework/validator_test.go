package framework

import (
	"github.com/go-playground/validator/v10"
	"testing"
)

type User struct {
	Mobile     string `json:"mobile" name:"手机号" validate:"required,len=11" example:"13800000000"`
	Source     string `json:"source" name:"注册渠道" validate:"required"`
	OpenId     string `json:"open_id" name:"微信open_id"`
	UnionId    string `json:"union_id" name:"微信union_id"`
	VerifyCode string `json:"verify_code" name:"短信验证码(微信登录传000000)" example:"1111" validate:"required"`
}

func  TestBackend_Validate(t *testing.T) {
	user := &User{}
	validate := validator.New()
	t.Log(validate.Struct(user))
}
