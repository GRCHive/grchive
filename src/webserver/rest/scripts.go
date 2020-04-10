package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type NewClientScriptInput struct {
	OrgId       int32  `json:"orgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func newClientScript(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewClientScriptInput{}
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

	script := core.ClientScript{
		OrgId:       inputs.OrgId,
		Name:        inputs.Name,
		Description: inputs.Description,
	}

	err = database.NewClientScript(&script, role)
	if err != nil {
		core.Warning("Failed to create new client script: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(script)
}

type UpdateClientScriptInput struct {
	OrgId       int32  `json:"orgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ScriptId    int64  `json:"scriptId"`
}

func updateClientScript(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UpdateClientScriptInput{}
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

	script := core.ClientScript{
		Id:          inputs.ScriptId,
		OrgId:       inputs.OrgId,
		Name:        inputs.Name,
		Description: inputs.Description,
	}

	err = database.UpdateClientScript(&script, role)
	if err != nil {
		core.Warning("Failed to create new client script: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(script)
}

type AllClientScriptsInput struct {
	OrgId int32 `webcore:"orgId"`
}

func allClientScripts(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllClientScriptsInput{}
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

	data, err := database.AllClientScriptsForOrganization(inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get all client data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(data)
}

type GetClientScriptInput struct {
	OrgId    int32 `webcore:"orgId"`
	ScriptId int64 `webcore:"scriptId"`
}

func getClientScript(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetClientScriptInput{}
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

	data, err := database.GetClientScriptFromId(inputs.ScriptId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get client data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(data)
}

type DeleteClientScriptInput struct {
	OrgId    int32 `json:"orgId"`
	ScriptId int64 `json:"scriptId"`
}

func deleteClientScript(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteClientScriptInput{}
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

	err = database.DeleteClientScript(inputs.ScriptId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete client data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
