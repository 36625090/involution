package logical

type Reply struct {
	Code       int         `json:"code" xml:"code" name:"业务状态码"`
	Data       interface{} `json:"data,omitempty" xml:"data" name:"数据体，任意数据"`
	Message    string      `json:"message,omitempty" xml:"message" name:"业务返回的消息"`
	Pagination *Pagination `json:"pagination,omitempty" xml:"pagination" name:"分页信息"`
}

type Pagination struct {
	Next       int `json:"next" name:"下一页"`
	Prev       int `json:"prev" name:"上一页"`
	Page       int `json:"page" name:"页码索引（从0开始）"`
	Size       int `json:"size" name:"分页大小"`
	Total      int `json:"total" name:"总条数"`
	TotalPages int `json:"totalPages" name:"总页数"`
}

// DocumentResponse
// 客户端请求schema的返回结构体
type DocumentResponse struct {
	Name      string    `json:"name"`
	Backend   string    `json:"backend"`
	Documents Documents `json:"namespaces"`
}
