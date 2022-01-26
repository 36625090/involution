package codes

type ReturnCode int

func (c ReturnCode) Int() int {
	return int(c)
}
func (c ReturnCode) Int64() int64 {
	return int64(c)
}

const (
	CodeSuccess                 ReturnCode = 0
	CodeFailure                 ReturnCode = 1
	CodeFatal                   ReturnCode = 2
	CodeServerInternalError     ReturnCode = 1001
	CodeDataValidateException   ReturnCode = 1002
	CodeRequestHeaderMissing    ReturnCode = 1003
	CodeEndpointNotFound        ReturnCode = 1004
	CodeOperationNotFound       ReturnCode = 1005
	CodeOperationHandlerIssue   ReturnCode = 1006
	CodeBackendIssue            ReturnCode = 1007
	CodeInvalidSignature        ReturnCode = 1008
	CodeInvalidRequestParameter ReturnCode = 1009
	CodeBindRequestData         ReturnCode = 1010

	CodeFailedDecodeArgs ReturnCode = 2001
	CodeServiceException ReturnCode = 3001
	CodeUnauthorized     ReturnCode = 4001
)
