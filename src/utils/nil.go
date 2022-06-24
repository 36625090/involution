package utils

import (
	"reflect"
)

func IsNil(v interface{})bool{
	if nil == v{return true}
	vType := reflect.ValueOf(v)
	switch vType.Kind() {
	case reflect.Ptr:
		if vType.Elem().Kind() == reflect.Struct{
			return IsNil(vType.Elem())
		}
		return vType.IsNil() || vType.IsZero()
	case reflect.Struct:
		return reflect.DeepEqual(v, struct{}{})
	case reflect.String:
		return v.(string) == ""
		default:
			return vType.Type() == nil
	}
}
