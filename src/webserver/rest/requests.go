package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type AllGenericRequestsInputs struct {
	OrgId int32 `webcore:"orgId"`
}

func allGenericRequestsScripts(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllGenericRequestsInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	reqs, err := database.GetGenericRequestsForScriptsInOrg(inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get requests: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(reqs)
}

type GetGenericRequestInputs struct {
	OrgId int32 `webcore:"orgId"`
}

func getGenericRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetGenericRequestScriptInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request, ok := r.Context().Value(webcore.GenericRequestContextKey).(*core.GenericRequest)
	if !ok || request == nil {
		core.Warning("Failed to get request from context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ret := struct {
		Request  *core.GenericRequest
		Approval *core.GenericApproval
	}{
		Request:  request,
		Approval: nil,
	}

	jsonWriter.Encode(ret)
}

type GetGenericRequestScriptInputs struct {
	OrgId int32 `webcore:"orgId"`
}

func getGenericRequestScript(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetGenericRequestScriptInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request, ok := r.Context().Value(webcore.GenericRequestContextKey).(*core.GenericRequest)
	if !ok || request == nil {
		core.Warning("Failed to get request from context")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code, err := database.GetCodeFromScriptRequestId(request.Id, role)
	if err != nil {
		core.Warning("Failed to get linked code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	script, err := database.GetScriptForCode(code.Id, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get linked script: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ret := struct {
		Script  *core.ClientScript
		Code    *core.ManagedCode
		OneTime core.NullTime
		RRule   core.NullString
		Params  map[string]interface{}
	}{
		Script: script,
		Code:   code,
	}

	ret.OneTime, ret.RRule, err = database.GetRunScheduleForScriptRequest(request.Id, role)
	if err != nil {
		core.Warning("Failed to get run request schedule: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ret.OneTime.NullTime.Valid || ret.RRule.NullString.Valid {
		ret.Params, err = database.GetParametersForScheduledScriptRunRequest(request.Id, role)
	} else {
		ret.Params, err = database.GetParametersForImmediateScriptRunRequest(request.Id, role)
	}

	if err != nil {
		core.Warning("Failed to get request run parameters: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(ret)
}
