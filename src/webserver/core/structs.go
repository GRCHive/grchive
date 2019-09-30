package core

import "reflect"

func StructToMap(inStruct interface{}) map[string]interface{} {
	retMap := make(map[string]interface{})

	rv := reflect.ValueOf(inStruct)
	for i := 0; i < rv.NumField(); i++ {
		retMap[rv.Type().Field(i).Name] = rv.Field(i).Interface()
	}

	return retMap
}
