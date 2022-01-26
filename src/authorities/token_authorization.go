package authorities

import (
	"context"
	"errors"
	"fmt"
)


type TokenAuthorization struct {
	authorization
	tokenHandler TokenHandler
}

func NewAuthorization(settings *Settings, handler TokenHandler) (Authorization,error) {
	if nil == settings{
		return nil, errors.New("authorization settings is nil")
	}

	a :=  &TokenAuthorization{}
	a.tokenHandler = handler
	a.settings = settings
	return a, nil
}

func (j *TokenAuthorization) Authentication(ctx context.Context, token string) (*Authorized, error) {
	authorized, err :=j.tokenHandler.ParseToken(token)
	if err != nil {
		return nil, fmt.Errorf("parse token: %v", err)
	}

	return authorized, nil
}
func (j *TokenAuthorization)TokenHandler()TokenHandler{
	return j.tokenHandler
}
