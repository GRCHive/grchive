package core

import (
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"time"
)

func FlexibleJsonStructUnmarshal(data []byte, out interface{}) error {
	rawMap := map[string]interface{}{}
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	outType := reflect.TypeOf(out)
	if outType.Kind() != reflect.Ptr {
		return errors.New("Can not unmarshal into non-pointer.")
	}

	outType = outType.Elem()
	outValue := reflect.ValueOf(out).Elem()

	for i := 0; i < outType.NumField(); i++ {
		fieldType := outType.Field(i)
		fieldValue := outValue.Field(i)

		validLookups := []string{
			fieldType.Name,
			fieldType.Tag.Get("db"),
		}

		jsonVal := GetFromMapWithKeyOptions(rawMap, validLookups...)

		var dataValue reflect.Value

		if jsonVal != nil {
			switch fieldType.Type {
			case Int64ReflectType:
				dataValue = reflect.ValueOf(int64(math.Round(jsonVal.(float64))))
			case Int32ReflectType:
				dataValue = reflect.ValueOf(int32(math.Round(jsonVal.(float64))))
			case StringReflectType:
				dataValue = reflect.ValueOf(jsonVal.(string))
			case NullInt64ReflectType:
				dataValue = reflect.ValueOf(CreateNullInt64(int64(math.Round(jsonVal.(float64)))))
			case NullTimeReflectType:
				t, err := time.Parse(time.RFC3339, jsonVal.(string))
				if err != nil {
					return err
				}
				dataValue = reflect.ValueOf(CreateNullTime(t))
			case TimeReflectType:
				t, err := time.Parse(time.RFC3339, jsonVal.(string))
				if err != nil {
					return err
				}
				dataValue = reflect.ValueOf(t)
			case BoolReflectType:
				dataValue = reflect.ValueOf(jsonVal.(bool))
			default:
				return errors.New("Flexible JSON unmarshal does not supported this type: " + fieldType.Type.Name())
			}
		} else {
			switch fieldType.Type {
			case NullInt64ReflectType:
				dataValue = reflect.ValueOf(NullInt64{})
			case NullTimeReflectType:
				dataValue = reflect.ValueOf(NullTime{})
			default:
				return errors.New("Nil value not handled for this type:" + fieldType.Type.Name())
			}
		}

		if !fieldValue.CanSet() {
			return errors.New("Can't set field: " + fieldType.Name)
		}
		fieldValue.Set(dataValue)
	}

	return nil
}
