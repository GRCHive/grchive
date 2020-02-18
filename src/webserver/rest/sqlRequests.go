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

type AllSqlRequestInput struct {
	DbId  core.NullInt64 `webcore:"dbId,optional"`
	OrgId int32          `webcore:"orgId"`
}

func allSqlRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllSqlRequestInput{}
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

	var data []*core.DbSqlQueryRequest
	if inputs.DbId.NullInt64.Valid {
		data, err = database.GetAllSqlRequestsForDb(inputs.DbId.NullInt64.Int64, inputs.OrgId, role)
	} else {
		data, err = database.GetAllSqlRequestsForOrg(inputs.OrgId, role)
	}

	if err != nil {
		core.Warning("Failed to get SQL requests: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(data)
}

type StatusSqlRequestInput struct {
	RequestId int64 `webcore:"requestId"`
	OrgId     int32 `webcore:"orgId"`
}

func statusSqlRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := StatusSqlRequestInput{}
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

	status, err := database.GetSqlRequestStatus(inputs.RequestId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get SQL request status: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(status)
}

type GetSqlRequestInput struct {
	RequestId int64 `webcore:"requestId"`
	OrgId     int32 `webcore:"orgId"`
}

func getSqlRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetSqlRequestInput{}
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

	request, err := database.GetSqlRequest(inputs.RequestId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get SQL request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	approval, err := database.GetSqlRequestStatus(inputs.RequestId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get SQL request status: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Request  *core.DbSqlQueryRequest
		Approval *core.DbSqlQueryRequestApproval
	}{
		Request:  request,
		Approval: approval,
	})
}

type UpdateSqlRequestInput struct {
	RequestId   int64  `json:"requestId"`
	OrgId       int32  `json:"orgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func updateSqlRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UpdateSqlRequestInput{}
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

	request := core.DbSqlQueryRequest{
		Id:          inputs.RequestId,
		OrgId:       inputs.OrgId,
		Name:        inputs.Name,
		Description: inputs.Description,
	}

	err = database.UpdateSqlQueryRequest(&request, role)
	if err != nil {
		core.Warning("Failed to create sql query request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(request)
}

type DeleteSqlRequestInput struct {
	RequestId int64 `json:"requestId"`
	OrgId     int32 `json:"orgId"`
}

func deleteSqlRequest(w http.ResponseWriter, r *http.Request) {
	inputs := UpdateSqlRequestInput{}
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

	err = database.DeleteSqlQueryRequest(inputs.RequestId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to create sql query request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type ModifyStatusSqlRequestInput struct {
	RequestId int64  `json:"requestId"`
	OrgId     int32  `json:"orgId"`
	Approve   bool   `json:"approve"`
	Reason    string `json:"reason"`
}

func modifyStatusSqlRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := ModifyStatusSqlRequestInput{}
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

	approval := core.DbSqlQueryRequestApproval{
		RequestId:        inputs.RequestId,
		OrgId:            inputs.OrgId,
		ResponseTime:     time.Now().UTC(),
		ResponsderUserId: userId,
		Response:         inputs.Approve,
		Reason:           inputs.Reason,
	}

	tx := database.CreateTx()

	err = database.UpdateRequestStatusWithTx(&approval, role, tx)
	if err != nil {
		tx.Rollback()
		core.Warning("Failed to update SQL request status: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if inputs.Approve {
		runCode, rawCode, err := webcore.GenerateRandomRunCode(approval.RequestId, approval.OrgId)
		if err != nil {
			tx.Rollback()
			core.Warning("Failed to generate random run code: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = database.CreateNewRunCodeWithTx(runCode, role, tx)
		if err != nil {
			tx.Rollback()
			core.Warning("Failed to create new run code in DB: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = webcore.SendRunCodeViaEmail(runCode, rawCode)
		if err != nil {
			tx.Rollback()
			core.Warning("Failed to send run code via email: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		core.Warning("Failed to commit SQL request status: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(approval)
}
