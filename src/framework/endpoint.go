package framework

import (
	"github.com/36625090/involution/logical"
	"reflect"
	"strings"
)
// EndpointAppend 生成endpoint数组
// list.
func EndpointAppend(paths ...[]*Endpoint) []*Endpoint {
	var result []*Endpoint
	for _, ps := range paths {
		result = append(result, ps...)
	}
	return result
}

type Endpoint struct {
	Pattern     string
	Description string
	Operations  map[string]OperationHandler
}

// OperationHandler operation接口
type OperationHandler interface {
	Handler() OperationFunc
	Properties() OperationProperties
}

// OperationProperties callback function操作
type OperationProperties struct {
	Description string
	Input       reflect.Type `json:"-"`
	Output      reflect.Type `json:"-"`
	Errors      logical.Errors     `json:"errors"`
}


// EndpointOperation is a concrete implementation of OperationHandler.
type EndpointOperation struct {
	Callback    OperationFunc
	Description string
	Input       reflect.Type
	Output      reflect.Type
	Errors      logical.Errors
}

func (p *EndpointOperation) Handler() OperationFunc {
	return p.Callback
}

func (p *EndpointOperation) Properties() OperationProperties {
	return OperationProperties{
		Description: strings.TrimSpace(p.Description),
		Input:       p.Input,
		Output:      p.Output,
		Errors:      p.Errors,
	}
}

