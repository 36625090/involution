package logical

import (
	"fmt"
	"github.com/36625090/involution/logical/codes"
)

type WrapperError struct {
	Code  codes.ReturnCode `json:"code"`
	Scope string           `json:"scope"`
	Err   error            `json:"err"`
}

func (w *WrapperError) Error() error {
	return w.Err
}

func (w *WrapperError) String() string {
	return fmt.Sprintf("WrapperError: code=%d scope=%s err=%s", w.Code, w.Scope, w.Err)
}
