package webcore

import (
	"errors"
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"net/http"
	"reflect"
	"strconv"
)

var boolReflectType = reflect.TypeOf((bool)(false))
var int64ReflectType = reflect.TypeOf((int64)(0))
var stringReflectType = reflect.TypeOf((string)(""))

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

		// TODO: Allow optional fields?
		requestParamName, ok := fieldType.Tag.Lookup("webcore")
		if !ok {
			requestParamName = fieldType.Name
		}

		data := r.Form[requestParamName]
		if len(data) == 0 {
			return errors.New("Could not find request param: " + requestParamName)
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
		case stringReflectType:
			dataValue = reflect.ValueOf(data[0])
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
