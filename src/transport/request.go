package transport

import (
	"sort"
	"strings"
)

//Request 客户端请求数据结构 Data会被解码传到后端真正的服务逻辑
type Request struct {
	Method    string `json:"method" binding:"required"`
	Data      string `json:"data" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
	Version   string `json:"version" binding:"required"`
	Sign      string `json:"sign" binding:"required"`
	SignType  string `json:"sign_type" choices:"md5" binding:"required"`
}

func (r *Request) Backend() string {
	return strings.Split(r.Method, ".")[0]
}

func (r *Request) Endpoint() string {
	return strings.Split(r.Method, ".")[1]
}

func (r *Request) Operation() string {
	return strings.Split(r.Method, ".")[2]
}

func (r *Request) Keys() []string {
	keys := []string{"method", "data", "timestamp", "version", "sign", "sign_type"}
	sort.Strings(keys)
	return keys
}

func (r *Request) Map() map[string]interface{} {
	params := make(map[string]interface{})
	params["data"] = r.Data
	params["method"] = r.Method
	params["sign_type"] = r.SignType
	params["timestamp"] = r.Timestamp
	params["version"] = r.Version
	params["sign"] = r.Sign
	return params
}
