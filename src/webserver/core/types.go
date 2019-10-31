package core

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

type NullInt64 struct {
	sql.NullInt64
}

func (v NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

var BoolReflectType = reflect.TypeOf((bool)(false))
var Int64ReflectType = reflect.TypeOf((int64)(0))
var NullInt64ReflectType = reflect.TypeOf(NullInt64{})
var Int32ReflectType = reflect.TypeOf((int32)(0))
var StringReflectType = reflect.TypeOf((string)(""))
var Int64ArrayReflectType = reflect.TypeOf(([]int64)([]int64{}))
