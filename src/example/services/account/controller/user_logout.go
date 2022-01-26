package controller

import (
	"context"
	"encoding/json"
	"github.com/36625090/involution/logical"
	"github.com/36625090/involution/logical/codes"
)

func (b *backend) userLogout(ctx context.Context, args *logical.Args, reply *logical.Reply) *logical.WrapperError {


	cli, err := b.LBAdapter.Client("userservice","").RestyClient()
	if err != nil {
		return &logical.WrapperError{
			Code: codes.CodeServiceException,
			Err:  err,
		}
	}

	var body = map[string]interface{}{
		"mobilePhone":"22222222222",
		"productCode":"12345",
	}

	request :=  cli.GetRequest()
	request.SetHeader("x-trace-id", args.GetTraceID())
	response, err := request.SetBody(body).
		SetHeader("Biz-ProductId","22222").
		Post("/users/loginByMobilePhone")
	if nil != err {
		return &logical.WrapperError{
			Code: codes.CodeServiceException,
			Err:  err,
		}
	}

	var body2 map[string]interface{}
	err = json.Unmarshal(response.Body(), &body2)
	if err != nil {
		return nil
	}

	reply.Code = response.StatusCode()
	reply.Data = body2
	return  nil
}
