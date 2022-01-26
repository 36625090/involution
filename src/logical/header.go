package logical

type HeaderKey string

func (h HeaderKey) String() string {
	return string(h)
}

const (
	HeaderTraceIDKey       HeaderKey = "X-Trace-ID"
	HeaderClientIDKey      HeaderKey = "X-Client-ID"
	HeaderApplicationKey   HeaderKey = "X-Application"
	HeaderAuthorizationKey HeaderKey = "Authorization"
)

