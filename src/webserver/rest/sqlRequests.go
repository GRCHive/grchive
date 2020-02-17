package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"time"
)

type NewSqlRequestInput struct {
	QueryId     int64  `json:"queryId"`
	OrgId       int32  `json:"orgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func newSqlRequest(w http.ResponseWriter, r *http.Request) {
	inputs := NewSqlRequestInput{}
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

	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	userId, err := webcore.GetUserIdFromApiRequestContext(r)
	if err != nil {
		core.Warning("Failed to obtain key user id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request := core.DbSqlQueryRequest{
		QueryId:      inputs.QueryId,
		UploadTime:   time.Now().UTC(),
		UploadUserId: userId,
		OrgId:        inputs.OrgId,
		Name:         inputs.Name,
		Description:  inputs.Description,
	}

	err = database.CreateNewSqlQueryRequest(&request, role)
	if err != nil {
		core.Warning("Failed to create sql query request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(request)
}
