package core

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func (v NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	} else {
		return json.Marshal(nil)
	}
}

func CreateNullTime(v time.Time) NullTime {
	return NullTime{sql.NullTime{v, true}}
}

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

func (v *NullInt64) UnmarshalJSON(b []byte) error {
	var val int64
	err := json.Unmarshal(b, &val)
	if err != nil {
		return nil
	}
	v.NullInt64.Int64 = val
	v.NullInt64.Valid = true
	return nil
}

type NullInt32 struct {
	sql.NullInt32
}

func (v NullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int32)
	} else {
		return json.Marshal(nil)
	}
}

func CreateNullInt32(v int32) NullInt32 {
	return NullInt32{sql.NullInt32{v, true}}
}

var BoolReflectType = reflect.TypeOf((bool)(false))
var IntReflectType = reflect.TypeOf((int)(0))
var Int64ReflectType = reflect.TypeOf((int64)(0))
var NullInt64ReflectType = reflect.TypeOf(NullInt64{})
var Int32ReflectType = reflect.TypeOf((int32)(0))
var StringReflectType = reflect.TypeOf((string)(""))
var Int64ArrayReflectType = reflect.TypeOf([]int64{})
var StringArrayReflectType = reflect.TypeOf([]string{})
var TimeReflectType = reflect.TypeOf(time.Time{})
