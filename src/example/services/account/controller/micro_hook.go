package controller

import (
	"github.com/hashicorp/go-hclog"
	"net/http"
)

type microHook struct {
	hclog.Logger
}

func (m *microHook) Trace(req *http.Request, res *http.Response, err error)  {
	m.Info("trace micro service", "req", req.URL, "res", res, "err",err)
}