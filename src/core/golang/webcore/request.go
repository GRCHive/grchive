package webcore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const ApiKeyRequestHeaderKey string = "ApiKey"

func GetRawClientAPIKeyFromRequest(r *http.Request) core.RawApiKey {
	return core.RawApiKey(r.Header.Get(ApiKeyRequestHeaderKey))
}

func GetRiskFromRequestUrl(r *http.Request, role *core.Role) (*core.Risk, error) {
	urlRouteVars := mux.Vars(r)
	riskIdStr, ok := urlRouteVars[core.DashboardOrgRiskQueryId]
	if !ok {
		return nil, errors.New("No risk in request URL")
	}

	riskId, err := strconv.ParseInt(riskIdStr, 10, 64)
	if err != nil {
		return nil, err
	}
	return database.FindRisk(riskId, role)
}

func GetControlFromRequestUrl(r *http.Request, role *core.Role) (*core.Control, error) {
	urlRouteVars := mux.Vars(r)
	controlIdStr, ok := urlRouteVars[core.DashboardOrgControlQueryId]
	if !ok {
		return nil, errors.New("No control in request URL")
	}

	controlId, err := strconv.ParseInt(controlIdStr, 10, 64)
	if err != nil {
		return nil, err
	}
	return database.FindControl(controlId, role)
}

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

func GetUserIdFromRequestUrl(r *http.Request) (int64, error) {
	urlRouteVars := mux.Vars(r)
	userIdStr, ok := urlRouteVars[core.DashboardUserQueryId]
	if !ok {
		return 0, errors.New("No user id in request URL")
	}
	return strconv.ParseInt(userIdStr, 10, 64)
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

func IsRequestMultipartForm(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	return strings.Contains(contentType, "multipart/form-data")
}

func IsJsonRequest(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	return strings.Contains(contentType, "application/json")
}

func UnmarshalRequestForm(r *http.Request, output interface{}) error {
	if IsRequestMultipartForm(r) {
		if err := r.ParseMultipartForm(MaxMultipartFormMemoryBytes); err != nil {
			return err
		}
	} else if IsJsonRequest(r) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, output)
		if err != nil {
			return err
		}

		return nil
	} else {
		if err := r.ParseForm(); err != nil {
			return err
		}
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
		case core.BoolReflectType:
			boolValue, err := strconv.ParseBool(data[0])
			if err != nil {
				return err
			}
			dataValue = reflect.ValueOf(boolValue)
			break
		case core.Int64ReflectType:
			intValue, err := strconv.ParseInt(data[0], 10, 64)
			if err != nil {
				return err
			}
			dataValue = reflect.ValueOf(intValue)
			break
		case core.IntReflectType:
			intValue, err := strconv.ParseInt(data[0], 10, 64)
			if err != nil {
				return err
			}
			dataValue = reflect.ValueOf(int(intValue))
			break
		case core.NullInt64ReflectType:
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
		case core.NullInt32ReflectType:
			intValue, err := strconv.ParseInt(data[0], 10, 32)
			if err != nil {
				return err
			}
			dataValue = reflect.ValueOf(core.NullInt32{
				sql.NullInt32{
					Int32: int32(intValue),
					Valid: true,
				},
			})
			break
		case core.NullBoolReflectType:
			boolValue, err := strconv.ParseBool(data[0])
			if err != nil {
				return err
			}
			dataValue = reflect.ValueOf(core.CreateNullBool(boolValue))
			break
		case core.Int32ReflectType:
			intValue, err := strconv.ParseInt(data[0], 10, 32)
			if err != nil {
				return err
			}
			dataValue = reflect.ValueOf(int32(intValue))
			break
		case core.StringReflectType:
			dataValue = reflect.ValueOf(data[0])
			break
		case core.StringArrayReflectType:
			dataValue = reflect.ValueOf(data[:])
			break
		case core.Int64ArrayReflectType:
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
		case core.TimeReflectType:
			inputDate, err := time.Parse(time.RFC3339, data[0])
			if err != nil {
				return err
			}
			dataValue = reflect.ValueOf(inputDate)
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

func GetUserIdFromApiRequestContext(r *http.Request) (int64, error) {
	key, err := FindApiKeyInContext(r.Context())
	if err != nil {
		return -1, err
	}
	return key.UserId, nil
}
