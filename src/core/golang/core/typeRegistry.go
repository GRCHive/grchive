package core

import "reflect"

var TypeRegistry = map[string]reflect.Type{
	UserReflectType.String():           UserReflectType,
	ControlReflectType.String():        ControlReflectType,
	DocRequestReflectType.String():     DocRequestReflectType,
	SqlRequestReflectType.String():     SqlRequestReflectType,
	DocMetadataReflectType.String():    DocMetadataReflectType,
	GenericRequestReflectType.String(): GenericRequestReflectType,
	"string":                           StringReflectType,
}

func GetBaseType(in interface{}) reflect.Type {
	typ := reflect.TypeOf(in)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}
