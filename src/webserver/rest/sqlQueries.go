package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"time"
)

type AllDatabaseQueryInputs struct {
	DbId  int64 `webcore:"dbId"`
	OrgId int32 `webcore:"orgId"`
}

func allDatabaseQuery(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllDatabaseQueryInputs{}
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

	metadata, err := database.GetAllSqlQueryMetadataForDb(inputs.DbId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get all query metadata: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(metadata)
}

type GetDatabaseQueryInput struct {
	MetadataId int64 `webcore:"metadataId"`
	OrgId      int32 `webcore:"orgId"`
}

func getDatabaseQuery(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetDatabaseQueryInput{}
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

	queries, err := database.GetAllSqlQueryVersionsForMetadata(inputs.MetadataId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get all metadata version: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(queries)
}

type NewDatabaseQueryInput struct {
	DbId         int64  `json:"dbId"`
	OrgId        int32  `json:"orgId"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	UploadUserId int64  `json:"uploadUserId"`
	Query        string `json:"query"`
}

func newDatabaseQuery(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewDatabaseQueryInput{}
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

	tx := database.CreateTx()

	metadata := core.DbSqlQueryMetadata{
		DbId:        inputs.DbId,
		OrgId:       inputs.OrgId,
		Name:        inputs.Name,
		Description: inputs.Description,
	}

	err = database.CreateSqlQueryMetadataWithTx(&metadata, role, tx)
	if err != nil {
		core.Warning("Failed to create metadata: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	query := core.DbSqlQuery{
		MetadataId:   metadata.Id,
		UploadTime:   time.Now().UTC(),
		UploadUserId: inputs.UploadUserId,
		OrgId:        inputs.OrgId,
		Query:        inputs.Query,
	}

	err = database.CreateSqlQueryWithTx(&query, role, tx)
	if err != nil {
		core.Warning("Failed to create query: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		core.Warning("Failed to commit new query: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Metadata *core.DbSqlQueryMetadata
		Query    *core.DbSqlQuery
	}{
		Metadata: &metadata,
		Query:    &query,
	})
}

type UpdateDatabaseQueryInput struct {
	OrgId      int32 `json:"orgId"`
	MetadataId int64 `json:"metadataId"`
	Metadata   *struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"metadata"`
	Query *struct {
		Query        string `json:"query"`
		UploadUserId int64  `json:"uploadUserId"`
	} `json:"query"`
}

func updateDatabaseQuery(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UpdateDatabaseQueryInput{}
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

	tx := database.CreateTx()

	var retMetadata *core.DbSqlQueryMetadata = nil
	if inputs.Metadata != nil {
		retMetadata = &core.DbSqlQueryMetadata{
			Id:          inputs.MetadataId,
			OrgId:       inputs.OrgId,
			Name:        inputs.Metadata.Name,
			Description: inputs.Metadata.Description,
		}
		err = database.UpdateSqlQueryMetadataWithTx(retMetadata, role, tx)
		if err != nil {
			core.Warning("Failed to update metadata: " + err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	var retQuery *core.DbSqlQuery = nil
	if inputs.Query != nil {
		retQuery = &core.DbSqlQuery{
			MetadataId:   inputs.MetadataId,
			UploadTime:   time.Now().UTC(),
			UploadUserId: inputs.Query.UploadUserId,
			OrgId:        inputs.OrgId,
			Query:        inputs.Query.Query,
		}

		err = database.CreateSqlQueryWithTx(retQuery, role, tx)
		if err != nil {
			core.Warning("Failed to create query version: " + err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		core.Warning("Failed to commit update query: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Metadata *core.DbSqlQueryMetadata
		Query    *core.DbSqlQuery
	}{
		Metadata: retMetadata,
		Query:    retQuery,
	})
}

type DeleteDatabaseQueryInput struct {
	OrgId      int32 `json:"orgId"`
	MetadataId int64 `json:"metadataId"`
}

func deleteDatabaseQuery(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteDatabaseQueryInput{}
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

	err = database.DeleteSqlQuery(inputs.MetadataId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete query: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
