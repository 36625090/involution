package transport

type Codec interface {
	Keys()[]string
	Map() map[string]interface{}
}