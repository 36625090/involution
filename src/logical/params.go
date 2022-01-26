package logical

//Encodable
//编码解码接口
type Encodable interface {
	Encode() []byte
	Decode(out interface{}) error
}