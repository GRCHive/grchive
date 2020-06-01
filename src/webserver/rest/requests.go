package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"time"
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

func allGenericRequestsShellScripts(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqs, err := database.GetGenericRequestsForShellScriptsInOrg(org.Id)
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

	approval, err := database.GetGenericApprovalForRequest(request.Id, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get approval for request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ret := struct {
		Request  *core.GenericRequest
		Approval *core.GenericApproval
	}{
		Request:  request,
		Approval: approval,
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

func getGenericRequestShell(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	request, err := webcore.FindGenericRequestInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get request from context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get organization from context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ret := struct {
		Shell      *core.ShellScript
		Version    *core.ShellScriptVersion
		VersionNum int32
		Servers    []*core.Server
	}{}

	ret.Shell, err = database.GetShellScriptFromRequest(request.Id)
	if err != nil {
		core.Warning("Failed to get shell script: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ret.Version, err = database.GetShellScriptVersionFromRequest(request.Id)
	if err != nil {
		core.Warning("Failed to get shell script version: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ret.Servers, err = database.GetShellScriptServersFromRequest(request.Id)
	if err != nil {
		core.Warning("Failed to get shell script servers: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Need to get all script versions so we can precisely know what index the returned
	// shell script version is - this is so we can give an accurate version number. This
	// might get slow if the number of versions increase...
	// Doing a O(n) search just because i'm assuming the number of versions will stay somewhat low.
	allVersions, err := database.AllShellScriptVersions(ret.Shell.Id, org.Id)
	if err != nil {
		core.Warning("Failed to get all script versions: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for idx, v := range allVersions {
		if v.Id == ret.Version.Id {
			ret.VersionNum = int32(len(allVersions) - idx)
			break
		}
	}

	jsonWriter.Encode(ret)
}

type EditGenericRequestInputs struct {
	OrgId   int32               `json:"orgId"`
	Request core.GenericRequest `json:"request"`
}

func editGenericRequest(w http.ResponseWriter, r *http.Request) {
	inputs := EditGenericRequestInputs{}
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

	err = database.EditGenericRequest(request.Id, inputs.OrgId, inputs.Request, role)
	if err != nil {
		core.Warning("Failed to edit request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type DeleteGenericRequestInputs struct {
	OrgId int32 `json:"orgId"`
}

func deleteGenericRequest(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteGenericRequestInputs{}
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

	err = database.DeleteGenericRequest(request.Id, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete generic request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type ApproveDenyGenericRequestInputs struct {
	OrgId   int32  `json:"orgId"`
	Approve bool   `json:"approve"`
	Reason  string `json:"reason"`
}

func approveDenyScriptRunRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := ApproveDenyGenericRequestInputs{}
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

	approval := core.GenericApproval{
		RequestId:       request.Id,
		ResponseTime:    time.Now().UTC(),
		ResponderUserId: role.UserId,
		Response:        inputs.Approve,
		Reason:          inputs.Reason,
	}

	err = database.InsertGenericApproval(&approval, role)
	if err != nil {
		core.Warning("Failed to insert generic approval: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if approval.Response {
		// Trigger attached script run job -it's either gonna be an immediate run or
		// a scheduled run so figure that out first before calling the right fn to
		// dispatch the job.
		typ, err := database.GetGenericRequestType(request.Id)
		if err != nil {
			core.Warning("Failed to get request type: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch typ {
		case core.KGenReqImmediateScript:
			runId, err := database.GetScriptRunIdLinkedToGenericRequest(request.Id)
			if err != nil {
				core.Warning("Failed to get run id linked to request: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = webcore.RunAuthorizedScriptImmediate(runId, approval)
		case core.KGenReqScheduledScript:
			taskId, err := database.GetTaskIdLinkedToGenericRequest(request.Id)
			if err != nil {
				core.Warning("Failed to get task id linked to request: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Modify task to have RunId and ApprovalId in its payload before
			// notifying the task manager to pick up the task.
			err = database.AddDataToTaskPayload(taskId, map[string]interface{}{
				"approvalId": approval.Id,
			})
			if err != nil {
				core.Warning("Failed to modify task payload: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = webcore.RunAuthorizedScriptScheduled(taskId, approval)
		default:
			core.Warning("Invalid request type")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err != nil {
			core.Warning("Failed to start script job: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	jsonWriter.Encode(approval)
}

func approveDenyShellRunRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := ApproveDenyGenericRequestInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request, err := webcore.FindGenericRequestInContext(r.Context())
	if err != nil {
		core.Warning("Can't find request in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.FindRoleInContext(r.Context())
	if err != nil {
		core.Warning("Can't find role in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	approval := core.GenericApproval{
		RequestId:       request.Id,
		ResponseTime:    time.Now().UTC(),
		ResponderUserId: role.UserId,
		Response:        inputs.Approve,
		Reason:          inputs.Reason,
	}

	err = database.InsertGenericApproval(&approval, role)
	if err != nil {
		core.Warning("Failed to insert generic approval: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if approval.Response {
		// TODO:
	}

	jsonWriter.Encode(approval)
}

type GetGenericApprovalInputs struct {
	OrgId int32 `webcore:"orgId"`
}

func getGenericApproval(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetGenericApprovalInputs{}
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

	approval, err := database.GetGenericApprovalForRequest(request.Id, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get approval: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(approval)
}
