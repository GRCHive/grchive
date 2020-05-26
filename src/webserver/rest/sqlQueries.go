package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/db_api/utility"
	"gitlab.com/grchive/grchive/proto/sqlQuery"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
	"google.golang.org/grpc"
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
	QueryId    int64 `webcore:"queryId"`
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

	var queries []*core.DbSqlQuery
	var desiredMetadataId int64
	if inputs.MetadataId != -1 {
		queries, err = database.GetAllSqlQueryVersionsForMetadata(inputs.MetadataId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get all metadata version: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		desiredMetadataId = inputs.MetadataId
	} else {
		q, err := database.GetSqlQueryFromId(inputs.QueryId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get query: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		queries = append(queries, q)

		desiredMetadataId = q.MetadataId
	}

	metadata, err := database.GetSqlMetadataFromId(desiredMetadataId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get metadata: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Queries  []*core.DbSqlQuery
		Metadata *core.DbSqlQueryMetadata
	}{
		Queries:  queries,
		Metadata: metadata,
	})
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

type RunDatabaseQueryInput struct {
	QueryId int64  `json:"queryId"`
	OrgId   int32  `json:"orgId"`
	RunCode string `json:"runCode"`
}

func runDatabaseQuery(w http.ResponseWriter, r *http.Request) {
	inputs := RunDatabaseQueryInput{}
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

	userId, err := webcore.GetUserIdFromApiRequestContext(r)
	if err != nil {
		core.Warning("Failed to obtain key user id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	codes, err := database.FindRunCodesForQueryForUser(inputs.QueryId, inputs.OrgId, userId, role)
	if err != nil {
		core.Warning("Failed to find existing codes: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var code *core.DbSqlQueryRunCode

	for _, c := range codes {
		if c.HashedCode == webcore.HashRunCode(inputs.RunCode, c.Salt) &&
			!core.IsPastTime(time.Now().UTC(), c.ExpirationTime, 10) {
			code = c
			break
		}
	}

	if code == nil {
		core.Warning("Failed to find a matching, unexpired code.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	// Ideally we have some sort of execute permission??
	if !role.Permissions.HasAccess(core.ResourceDbSqlQuery, core.AccessView) {
		core.Warning("No permission to execute query.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := webcore.CreateGRPCClientConnection(core.EnvConfig.Grpc.QueryRunner, core.EnvConfig.Tls, grpc.WithTimeout(time.Second*30))
	if err != nil {
		core.Warning("Failed to connect to GRPC: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer conn.Close()

	client := sqlQuery.NewQueryRunnerClient(conn)

	// 30 seconds is probably sufficient...this probably needs to be configurable at some point though.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	transitKey := fmt.Sprintf("sqlquery-%d", inputs.QueryId)

	resp, err := client.RunSqlQuery(ctx, &sqlQuery.SqlRunnerRequest{
		VaultResultPath: transitKey,
		QueryId:         inputs.QueryId,
		OrgId:           inputs.OrgId,
	})

	if err != nil {
		core.Warning("Failed to run query: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	decrypt, err := vault.TransitDecrypt(transitKey, resp.EncryptedData)
	if err != nil {
		core.Warning(fmt.Sprintf("Failed to decrypt result: %s [Error: %s]", string(resp.EncryptedData), err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Decrypt is either an utility.SqlQueryResult or a string depending on whether
	// Success is true or false.
	if !resp.Success {
		jsonWriter.Encode(struct {
			Data    string
			Success bool
		}{
			Data:    string(decrypt),
			Success: false,
		})
	} else {
		data := struct {
			Data    *utility.SqlQueryResult
			Success bool
		}{
			Data:    nil,
			Success: true,
		}

		err = json.Unmarshal(decrypt, &data.Data)
		if err != nil {
			core.Warning("Failed to parse result: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonWriter.Encode(data)
	}

	// We should able to ignore this error.	Worst case scenario it'll just expire.
	err = database.MarkRunCodeAsUsed(code.HashedCode, code.RequestId, code.OrgId, core.ServerRole)
	if err != nil {
		core.Warning("Failed to mark used: " + err.Error())
	}
}
