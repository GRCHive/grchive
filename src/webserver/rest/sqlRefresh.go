package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type AllDatabaseRefreshInputs struct {
	DbId  int64 `webcore:"dbId"`
	OrgId int32 `webcore:"orgId"`
}

func allDatabaseRefresh(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllDatabaseRefreshInputs{}
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

	refresh, err := database.GetAllDatabaseRefresh(inputs.DbId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get refresh: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(refresh)
}

type GetDatabaseRefreshInputs struct {
	RefreshId int64 `webcore:"refreshId"`
	OrgId     int32 `webcore:"orgId"`
}

func getDatabaseRefresh(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetDatabaseRefreshInputs{}
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

	refresh, err := database.GetDatabaseRefresh(inputs.RefreshId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get refresh: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jsonWriter.Encode(refresh)
}

type NewRefreshInput struct {
	DbId  int64 `json:"dbId"`
	OrgId int32 `json:"orgId"`
}

func newDatabaseRefresh(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewRefreshInput{}
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

	refresh, err := database.CreateNewDatabaseRefresh(inputs.DbId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to create DB Refresh: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jsonWriter.Encode(refresh)
}
