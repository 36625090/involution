package framework

import (
	"fmt"
	"github.com/36625090/involution/logical"
)

func (b *Backend) initDocumentsOnce() {
	_ = b.initDocuments()
}

func (b *Backend) initDocuments()error {
	endpoints := logical.Documents{}
	for _, ns := range b.Endpoints {
		if ns.Description == "" {
			return fmt.Errorf("endpoint[%s] description required", ns.Pattern)
		}
		endpoint := logical.Document{
			Endpoint:    ns.Pattern,
			Description: ns.Description,
			Operations:  make(map[string]*logical.Operation),
		}
		for opt, handler := range ns.Operations {
			properties := handler.Properties()
			if properties.Description == "" {
				return descriptionError(ns.Pattern, opt)
			}
			input, err := logical.Fields(properties.Input)
			if err != nil {
				return endpointError(ns.Pattern, opt, err)
			}
			output, err := logical.Fields(properties.Output)
			if err != nil {
				return endpointError(ns.Pattern, opt, err)
			}
			operation := &logical.Operation{
				Description: properties.Description,
				Input:       input,
				Output:      output,
				Errors:      properties.Errors,
			}
			endpoint.Operations[opt] = operation
		}
		endpoints = append(endpoints, &endpoint)
	}
	b.documents = endpoints
	return nil
}

func endpointError(pattern string, operation string, err error) error {
	return fmt.Errorf("endpoint[%s] operation[%s] %s", pattern, operation, err)
}
func descriptionError(pattern string, operation string) error {
	return fmt.Errorf("endpoint[%s] operation[%s] Description required", pattern, operation)
}
