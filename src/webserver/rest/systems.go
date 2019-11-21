package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type NewSystemInputs struct {
	OrgId       int32  `json:"orgId"`
	Name        string `json:"name"`
	Purpose     string `json:"purpose"`
	Description string `json:"description"`
}

type SystemAllInputs struct {
	OrgId int32 `webcore:"orgId"`
}

type GetSystemInputs struct {
	SysId int64 `webcore:"sysId"`
	OrgId int32 `webcore:"orgId"`
}

type EditSystemInputs struct {
	SysId       int64  `json:"sysId"`
	OrgId       int32  `json:"orgId"`
	Name        string `json:"name"`
	Purpose     string `json:"purpose"`
	Description string `json:"description"`
}

type DeleteSystemInputs struct {
	SysId int64 `json:"sysId"`
	OrgId int32 `json:"orgId"`
}

func newSystem(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewSystemInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sys := core.System{
		OrgId:       inputs.OrgId,
		Name:        inputs.Name,
		Purpose:     inputs.Purpose,
		Description: inputs.Description,
	}

	err = database.CreateNewSystem(&sys, role)
	if err != nil {
		core.Warning("Failed to create new system: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(sys)
}

func getAllSystems(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := SystemAllInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	systems, err := database.GetAllSystemsForOrg(inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get all systems: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(systems)
}

func getSystem(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetSystemInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sys, err := database.GetSystem(inputs.SysId, org.Id, role)
	if err != nil {
		core.Warning("Failed to get system: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		System *core.System
	}{
		System: sys,
	})
}

func editSystem(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := EditSystemInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sys := core.System{
		Id:          inputs.SysId,
		OrgId:       inputs.OrgId,
		Name:        inputs.Name,
		Purpose:     inputs.Purpose,
		Description: inputs.Description,
	}

	err = database.EditSystem(&sys, role)
	if err != nil {
		core.Warning("Failed to edit system: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(sys)
}

func deleteSystem(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteSystemInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = database.DeleteSystem(inputs.SysId, org.Id, role)
	if err != nil {
		core.Warning("Failed to delete system: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
