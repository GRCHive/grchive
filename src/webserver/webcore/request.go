package webcore

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var boolReflectType = reflect.TypeOf((bool)(false))
var int64ReflectType = reflect.TypeOf((int64)(0))
var nullInt64ReflectType = reflect.TypeOf(core.NullInt64{})
var int32ReflectType = reflect.TypeOf((int32)(0))
var stringReflectType = reflect.TypeOf((string)(""))
var int64ArrayReflectType = reflect.TypeOf(([]int64)([]int64{}))

func GetOrganizationFromRequestUrl(r *http.Request) (*core.Organization, error) {
	urlRouteVars := mux.Vars(r)
	orgGroupName, ok := urlRouteVars[core.DashboardOrgOrgQueryId]
	if !ok {
		return nil, errors.New("No valid organization in request URL")
	}

	org, err := database.FindOrganizationFromGroupName(orgGroupName)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func GetUserEmailFromRequestUrl(r *http.Request) (string, error) {
	urlRouteVars := mux.Vars(r)
	email, ok := urlRouteVars[core.DashboardUserQueryId]
	if !ok {
		return "", errors.New("No email in request URL")
	}

	return email, nil
}

func GetProcessFlowIdFromRequest(r *http.Request) (int64, error) {
	urlRouteVars := mux.Vars(r)
	id, ok := urlRouteVars[core.ProcessFlowQueryId]
	if !ok {
		return 0, errors.New("No process flow id in request URL")
	}

	val, err := strconv.ParseInt(id, 10, 64)
	return val, err
}

func UnmarshalRequestForm(r *http.Request, output interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	interfaceType := reflect.TypeOf(output).Elem()
	interfaceValue := reflect.ValueOf(output).Elem()

	for i := 0; i < interfaceType.NumField(); i++ {
		fieldType := interfaceType.Field(i)
		fieldValue := interfaceValue.Field(i)

		var requestParamName string
		var optional bool

		tag, ok := fieldType.Tag.Lookup("webcore")
		if !ok {
			requestParamName = fieldType.Name
			optional = false
		} else {
			splitTag := strings.Split(tag, ",")
			requestParamName = splitTag[0]
			// This probably needs to be more flexible instead of just assuming
			// that the optional tag will come as the 2nd parameter?
			optional = len(splitTag) > 1 && splitTag[1] == "optional"
		}

		data := r.Form[requestParamName]
		if len(data) == 0 {
			if !optional {
				return errors.New("Could not find request param: " + requestParamName)
			} else {
				continue
			}
		}

		var dataValue reflect.Value
		// Convert the string to the proper type.
		// Only handle the case for now where a single parameter value is passed in.
		switch fieldType.Type {
		case boolReflectType:
			boolValue, err := strconv.ParseBool(data[0])
			if err != nil {
				return err
			}
			dataValue = reflect.ValueOf(boolValue)
			break
		case int64ReflectType:
			intValue, err := strconv.ParseInt(data[0], 10, 64)
			if err != nil {
				return err
			}
			dataValue = reflect.ValueOf(intValue)
			break
		case nullInt64ReflectType:
			intValue, err := strconv.ParseInt(data[0], 10, 64)
			if err != nil {
				return err
			}
			dataValue = reflect.ValueOf(core.NullInt64{
				sql.NullInt64{
					Int64: intValue,
					Valid: true,
				},
			})
			break
		case int32ReflectType:
			intValue, err := strconv.ParseInt(data[0], 10, 32)
			if err != nil {
				return err
			}
			dataValue = reflect.ValueOf(int32(intValue))
			break
		case stringReflectType:
			dataValue = reflect.ValueOf(data[0])
			break
		case int64ArrayReflectType:
			arr := make([]int64, len(data))
			for idx, val := range data {
				intValue, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return err
				}
				arr[idx] = int64(intValue)
			}
			dataValue = reflect.ValueOf(arr)
			break
		default:
			return errors.New("Unsupported type: " + fieldType.Name)
		}

		if !fieldValue.CanSet() {
			return errors.New("Can't set field: " + fieldType.Name)
		}
		fieldValue.Set(dataValue)
	}
	return nil
}
