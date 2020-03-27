package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type NewClientDataInput struct {
	OrgId        int32                  `json:"orgId"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	SourceId     int64                  `json:"sourceId"`
	SourceTarget map[string]interface{} `json:"sourceTarget"`
}

func newClientData(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewClientDataInput{}
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

	data := core.ClientData{
		OrgId:       inputs.OrgId,
		Name:        inputs.Name,
		Description: inputs.Description,
	}

	tx, err := database.CreateAuditTrailTx(role)
	if err != nil {
		core.Warning("Failed to create tx: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.NewClientDataWithTx(&data, role, tx)
	if err != nil {
		tx.Rollback()
		core.Warning("Failed to create client data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.LinkClientDataToSourceWithTx(data.Id, inputs.SourceId, inputs.SourceTarget, inputs.OrgId, role, tx)
	if err != nil {
		tx.Rollback()
		core.Warning("Failed to link client data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		core.Warning("Failed to commit: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(core.FullClientDataWithLink{
		Data: data,
		Link: core.DataSourceLink{
			OrgId:        inputs.OrgId,
			DataId:       data.Id,
			SourceId:     inputs.SourceId,
			SourceTarget: inputs.SourceTarget,
		},
	})
}

type AllClientDataInput struct {
	OrgId int32 `webcore:"orgId"`
}

func allClientData(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllClientDataInput{}
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

	data, err := database.AllClientDataForOrganization(inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get all client data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(data)
}

type DeleteClientDataInput struct {
	OrgId  int32 `json:"orgId"`
	DataId int64 `json:"dataId"`
}

func deleteClientData(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteClientDataInput{}
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

	err = database.DeleteClientData(inputs.DataId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete client data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func allDataSourceOptions(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	// The only thing we can do for metadata requests is to make sure
	// a valid API key is found.
	_, err := webcore.FindApiKeyInContext(r.Context())
	if err != nil {
		core.Warning("Can't find API key: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	options, err := database.AllDataSourceOptions()
	if err != nil {
		core.Warning("Failed to get all data source options: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(options)
}
