package logical
type HttpMethod string
const (
	HttpMethodGET     HttpMethod = "GET"
	HttpMethodPOST    HttpMethod = "POST"
	HttpMethodPUT     HttpMethod = "PUT"
	HttpMethodDELETE  HttpMethod= "DELETE"
	HttpMethodPATCH   HttpMethod = "PATCH"
	HttpMethodOPTIONS HttpMethod = "OPTIONS"
	HttpMethodHEAD    HttpMethod = "HEAD"
	HttpMethodTRACE   HttpMethod = "TRACE"
	HttpMethodCONNECT HttpMethod = "CONNECT"
)

