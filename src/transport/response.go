package transport

import (
	"encoding/json"
	"sort"
)

// Response 给客户端的返回数据结构
type Response struct {
	Code      int                 `xml:"code" json:"code"`
	Message   string              `xml:"message" json:"message"`
	Content   interface{}         `xml:"content" json:"content"`
	TraceID   string              `json:"trace_id" xml:"trace_id"`
	Timestamp int64               `json:"timestamp" xml:"timestamp"`
	Sign      string              `json:"sign" xml:"sign"`
}

func (r *Response) Keys() []string {
	keys := []string{"code", "message", "content", "trace_id", "timestamp", "sign"}
	sort.Strings(keys)
	return keys
}

func (r *Response) Map() map[string]interface{} {
	bs, _ := json.Marshal(r.Content)
	params := make(map[string]interface{})
	params["code"] = r.Code
	params["message"] = r.Message
	params["content"] = string(bs)
	params["trace_id"] = r.TraceID
	params["timestamp"] = r.Timestamp
	params["sign"] = r.Sign
	return params
}