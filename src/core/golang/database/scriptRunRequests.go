package database

import (
	"errors"
	"github.com/jmoiron/sqlx/types"
	"gitlab.com/grchive/grchive/core"
	"reflect"
)

func GetRunScheduleForScriptRequest(requestId int64, role *core.Role) (core.NullTime, core.NullString, error) {
	retOneTime := core.NullTime{}
	retRRule := core.NullString{}

	rows, err := dbConn.Queryx(`
		SELECT ot.event_date_time, rt.rrule
		FROM generic_requests AS req
		INNER JOIN request_to_scheduled_task_link AS stl
			ON stl.request_id = req.id
		LEFT JOIN one_time_tasks AS ot
			ON ot.event_id = stl.task_id
		LEFT JOIN recurring_tasks AS rt
			ON rt.event_id = stl.task_id
		WHERE req.id = $1
	`, requestId)
	if err != nil {
		return retOneTime, retRRule, err
	}

	defer rows.Close()
	if !rows.Next() {
		return retOneTime, retRRule, nil
	}

	err = rows.Scan(&retOneTime, &retRRule)
	if err != nil {
		return retOneTime, retRRule, err
	}

	return retOneTime, retRRule, nil
}

func GetParametersForScheduledScriptRunRequest(requestId int64, role *core.Role) (map[string]interface{}, error) {
	rows, err := dbConn.Queryx(`
		SELECT task.task_data #> '{Payload, params}'
		FROM scheduled_tasks AS task
		INNER JOIN request_to_scheduled_task_link AS link
			ON link.task_id = task.id
		WHERE link.request_id = $1
	`, requestId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rows.Next()

	jsonData := types.JSONText{}
	err = rows.Scan(&jsonData)
	if err != nil {
		return nil, err
	}

	retMap := map[string]interface{}{}
	err = jsonData.Unmarshal(&retMap)
	if err != nil {
		return nil, err
	}

	return retMap, nil
}

func GetParametersForImmediateScriptRunRequest(requestId int64, role *core.Role) (map[string]interface{}, error) {
	rows, err := dbConn.Queryx(`
		SELECT p.run_id, p.param_name, p.vals, typ.golang_type
		FROM script_run_parameters AS p
		INNER JOIN request_to_script_run_link AS link
			ON link.run_id = p.run_id
		INNER JOIN client_script_code_parameters AS params
			ON params.name = p.param_name
		INNER JOIN supported_code_parameter_types AS typ
			ON typ.id = params.param_type
		WHERE link.request_id = $1
	`, requestId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	retMap := map[string]interface{}{}
	for rows.Next() {
		param := core.ScriptRunParameter{}
		jsonData := types.JSONText{}
		golangType := ""

		err = rows.Scan(&param.RunId, &param.ParamName, &jsonData, &golangType)
		if err != nil {
			return nil, err
		}

		reflectType, ok := core.TypeRegistry[golangType]
		if !ok {
			return nil, errors.New("Unsupported type: " + golangType)
		}

		val := reflect.New(reflectType).Interface()
		err = jsonData.Unmarshal(val)
		if err != nil {
			return nil, err
		}

		retMap[param.ParamName] = val
	}
	return retMap, nil
}
