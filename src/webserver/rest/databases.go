package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type NewDatabaseInputs struct {
	Name      string `json:"name"`
	OrgId     int32  `json:"orgId"`
	TypeId    int32  `json:"typeId"`
	OtherType string `json:"otherType"`
	Version   string `json:"version"`
}

type GetAllDatabaseInputs struct {
	OrgId int32 `webcore:"orgId"`
}

type GetDatabaseInputs struct {
	DbId  int64 `webcore:"dbId"`
	OrgId int32 `webcore:"orgId"`
}

type EditDatabaseInputs struct {
	DbId      int64  `json:"dbId"`
	Name      string `json:"name"`
	OrgId     int32  `json:"orgId"`
	TypeId    int32  `json:"typeId"`
	OtherType string `json:"otherType"`
	Version   string `json:"version"`
}

type DeleteDatabaseInputs struct {
	DbId  int64 `json:"dbId"`
	OrgId int32 `json:"orgId"`
}

func newDb(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewDatabaseInputs{}
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

	db := core.Database{
		Name:      inputs.Name,
		OrgId:     inputs.OrgId,
		TypeId:    inputs.TypeId,
		OtherType: inputs.OtherType,
		Version:   inputs.Version,
	}

	err = database.InsertNewDatabase(&db, role)
	if err != nil {
		core.Warning("Can't insert new database: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(db)
}

func getAllDb(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetAllDatabaseInputs{}
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

	dbs, err := database.GetAllDatabasesForOrg(org.Id, role)
	if err != nil {
		core.Warning("Can't get all databases: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(dbs)
}

func getDbTypes(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	apiKey, err := webcore.GetAPIKeyFromRequest(r)
	if apiKey == nil || err != nil {
		core.Warning("No API Key: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	types, err := database.GetAllSupportedDatabaseTypes(core.ServerRole)
	if err != nil {
		core.Warning("Can't get types: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(types)
}

func getDb(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetDatabaseInputs{}
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

	db, err := database.GetDb(inputs.DbId, org.Id, role)
	if err != nil {
		core.Warning("Can't get database: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Database *core.Database
	}{
		Database: db,
	})
}

func editDb(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := EditDatabaseInputs{}
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

	db := core.Database{
		Id:        inputs.DbId,
		Name:      inputs.Name,
		OrgId:     inputs.OrgId,
		TypeId:    inputs.TypeId,
		OtherType: inputs.OtherType,
		Version:   inputs.Version,
	}

	err = database.EditDb(&db, role)
	if err != nil {
		core.Warning("Can't edit database: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(db)
}

func deleteDb(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteDatabaseInputs{}
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

	err = database.DeleteDb(inputs.DbId, org.Id, role)
	if err != nil {
		core.Warning("Can't delete database: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
