package transport

import (
	"encoding/json"
	"github.com/36625090/involution/utils"
	"sort"
)

// Response 给客户端的返回数据结构
type Response struct {
	Code       int         `json:"code" xml:"code"`
	Message    string      `json:"message" xml:"message"`
	Content    interface{} `json:"content" xml:"content"`
	Pagination interface{} `json:"pagination" xml:"pagination"`
	TraceID    string      `json:"trace_id" xml:"trace_id"`
	Timestamp  int64       `json:"timestamp" xml:"timestamp"`
	Sign       string      `json:"sign" xml:"sign"`
}

func (r *Response) Keys() []string {
	keys := []string{"code", "message", "content", "trace_id", "timestamp", "sign","pagination"}
	sort.Strings(keys)
	return keys
}

type AX struct {

}
func (r *Response) Map() map[string]interface{} {

	params := make(map[string]interface{})
	params["code"] = r.Code
	params["message"] = r.Message
	params["trace_id"] = r.TraceID
	params["timestamp"] = r.Timestamp
	params["sign"] = r.Sign
	if ! utils.IsNil(r.Content){
		bs, _ := json.Marshal(r.Content)
		params["content"] = string(bs)
	}
	if ! utils.IsNil(r.Pagination){
		bs, _ := json.Marshal(r.Pagination)
		params["pagination"] = string(bs)
	}
	return params
}
