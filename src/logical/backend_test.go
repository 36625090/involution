package logical

import (
	"testing"
	"time"
)

func TestAuthorized_Encode(t *testing.T) {
	c := &BackendContext{
		Application:          "Test-Configuration",
	}

	go func(cfg *BackendContext) {
		for{
			select {
			case <-time.After(time.Second*3):
				t.Log(cfg.Application)
			}
		}
	}(c)

	go func(cfg *BackendContext) {
		for{
			select {
			case <-time.After(time.Second*4):
				cfg.Application = time.Now().String()
			}
		}
	}(c)
	go func(cfg *BackendContext) {
		for{
			select {
			case <-time.After(time.Second*4):
				cfg.Application = "I am "+time.Now().String()
			}
		}
	}(c)

	select {

	}
}