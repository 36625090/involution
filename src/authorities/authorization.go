package authorities

import (
	"context"
)

type Authorization interface {
	Settings() *Settings
	Authentication(ctx context.Context, token string) (*Authorized, error)
	TokenHandler()TokenHandler
}

type authorization struct {
	settings *Settings
}

func (j *authorization) Settings() *Settings {
	return j.settings
}