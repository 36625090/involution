package logical

import (
	"reflect"
	"strings"
)

type EmptyDocuments struct{}
type Errors map[int]string

type DocumentsReply struct {
	Documents Documents `json:"documents"`
}

type Documents []*Document
type Document struct {
	Endpoint    string                `json:"namespace,omitempty"`
	Description string                `json:"description"`
	Operations  map[string]*Operation `json:"operations"`
}

//Operation
//属性
type Operation struct {
	Description string   `json:"description"`
	Input       []*Field `json:"input,omitempty"`
	Output      []*Field `json:"output,omitempty"`
	Errors      Errors   `json:"errors,omitempty"`
}

//Field
//接口属性列
type Field struct {
	Field     string   `json:"field"`
	Name      string   `json:"name"`
	Kind      string   `json:"kind"`
	Required  bool     `json:"required"`
	Reference []*Field `json:"reference"`
	Example   string   `json:"example"`
	IsList    bool     `json:"is_list"`
}


func Fields(t reflect.Type) ([]*Field, error) {
	return getFields(t), nil
}

func getType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Ptr:
		return t.Elem()
	case reflect.Struct:
		return t
	case reflect.Slice:
		fallthrough
	case reflect.Map:
		return getType(t.Elem())
	default:
		return t
	}
}

func getFields(Type reflect.Type) []*Field {
	defer func() {
		recover()
	}()
	var fields []*Field
	for i := 0; i < Type.NumField(); i++ {
		field := new(Field)
		f := Type.Field(i)
		isList := f.Type.Kind() == reflect.Slice
		realType := getType(f.Type)
		kindString := realType.Kind().String()
		if realType.Kind() == reflect.Struct && realType.Name() != "Time" {
			field.Reference = getFields(realType)
		}
		if realType.Name() == "Time" {
			kindString = "datetime"
		}

		fValue := f.Tag.Get("json")
		if fValue == "" {
			fValue = f.Name
		}
		fName := f.Tag.Get("name")
		if fName == "" {
			fName = f.Name
		}
		example := f.Tag.Get("example")
		validate := f.Tag.Get("validate")
		required := validate != "" && strings.Contains(strings.ToLower(validate), "required")

		field.Name = fName
		field.Field = fValue
		field.Required = required
		field.Kind = kindString
		field.Example = example
		field.IsList = isList
		fields = append(fields, field)
	}
	return fields
}
