package core

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"
)

type NullString struct {
	sql.NullString
}

func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullString) UnmarshalJSON(b []byte) error {
	var val string
	err := json.Unmarshal(b, &val)
	if err != nil || string(b) == "null" {
		return nil
	}
	v.NullString.String = val
	v.NullString.Valid = true
	return nil
}

func CreateNullString(v string) NullString {
	return NullString{sql.NullString{v, true}}
}

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

func (v *NullTime) UnmarshalJSON(b []byte) error {
	var val time.Time
	err := json.Unmarshal(b, &val)
	if err != nil || string(b) == "null" {
		return nil
	}
	v.NullTime.Time = val
	v.NullTime.Valid = true
	return nil
}

func CreateNullTime(v time.Time) NullTime {
	return NullTime{sql.NullTime{v, true}}
}

func (a NullTime) Equal(b NullTime) bool {
	if a.NullTime.Valid != b.NullTime.Valid {
		return false
	}

	if !a.NullTime.Valid {
		return true
	}

	return a.NullTime.Time.Equal(b.NullTime.Time)
}

type NullInt64 struct {
	sql.NullInt64
}

func CreateNullInt64(v int64) NullInt64 {
	return NullInt64{sql.NullInt64{v, true}}
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
	if err != nil || string(b) == "null" {
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

func (v *NullInt32) UnmarshalJSON(b []byte) error {
	var val int32
	err := json.Unmarshal(b, &val)
	if err != nil || string(b) == "null" {
		return nil
	}
	v.NullInt32.Int32 = val
	v.NullInt32.Valid = true
	return nil
}

func CreateNullInt32(v int32) NullInt32 {
	return NullInt32{sql.NullInt32{v, true}}
}

type NullBool struct {
	sql.NullBool
}

func (v NullBool) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Bool)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullBool) UnmarshalJSON(b []byte) error {
	var val bool
	err := json.Unmarshal(b, &val)
	if err != nil || string(b) == "null" {
		return nil
	}
	v.NullBool.Bool = val
	v.NullBool.Valid = true
	return nil
}

func CreateNullBool(v bool) NullBool {
	return NullBool{sql.NullBool{v, true}}
}

var BoolReflectType = reflect.TypeOf((bool)(false))
var IntReflectType = reflect.TypeOf((int)(0))
var Int64ReflectType = reflect.TypeOf((int64)(0))
var NullInt64ReflectType = reflect.TypeOf(NullInt64{})
var NullInt32ReflectType = reflect.TypeOf(NullInt32{})
var NullBoolReflectType = reflect.TypeOf(NullBool{})
var NullStringReflectType = reflect.TypeOf(NullString{})
var NullTimeReflectType = reflect.TypeOf(NullTime{})
var Int32ReflectType = reflect.TypeOf((int32)(0))
var StringReflectType = reflect.TypeOf((string)(""))
var Int64ArrayReflectType = reflect.TypeOf([]int64{})
var StringArrayReflectType = reflect.TypeOf([]string{})
var TimeReflectType = reflect.TypeOf(time.Time{})

var RiskFilterDataType = reflect.TypeOf(RiskFilterData{})
var ControlFilterDataType = reflect.TypeOf(ControlFilterData{})
var DatabaseFilterDataType = reflect.TypeOf(DatabaseFilterData{})
var AuditTrailFilterDataType = reflect.TypeOf(AuditTrailFilterData{})

var UserReflectType = reflect.TypeOf(User{})
var ControlReflectType = reflect.TypeOf(Control{})
var DocRequestReflectType = reflect.TypeOf(DocumentRequest{})
var SqlRequestReflectType = reflect.TypeOf(DbSqlQueryRequest{})
var DocMetadataReflectType = reflect.TypeOf(ControlDocumentationFile{})

var DataSourceLinkReflectType = reflect.TypeOf(DataSourceLink{})
