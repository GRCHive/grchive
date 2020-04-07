package rest

import (
	"encoding/json"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"time"
)

type SaveCodeInput struct {
	OrgId  int32          `json:"orgId"`
	Code   string         `json:"code"`
	DataId core.NullInt64 `json:"dataId"`
}

func saveCode(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := SaveCodeInput{}
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

	if !inputs.DataId.NullInt64.Valid {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	managedCode := core.ManagedCode{
		OrgId:      inputs.OrgId,
		ActionTime: time.Now().UTC(),
	}

	// Determine GitPath from whether this is a data or script.
	if inputs.DataId.NullInt64.Valid {
		clientData, err := database.GetClientDataFromId(inputs.DataId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get client data: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// For now, assume Kotlin always. If we want to support more in the future we'll have to
		// somehow get this information from the user or something.
		managedCode.GitPath = fmt.Sprintf("src/main/kotlin/data/%s", clientData.Data.Filename("kt"))
	}

	// There is a possibility here that the link will fail after the storing to Gitea succeeds.
	// Do we care in that case? We can probably survive just losing a link since storing an
	// extra commit in Gitea won't hurt us.
	err = webcore.StoreManagedCodeToGitea(&managedCode, inputs.Code, role)
	if err != nil {
		core.Warning("Failed to store managed code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if inputs.DataId.NullInt64.Valid {
		err = database.LinkCodeToData(managedCode.Id, inputs.DataId.NullInt64.Int64, inputs.OrgId, role)
	}

	if err != nil {
		core.Warning("Failed to link code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(managedCode)
}

type GetCodeInput struct {
	OrgId  int32          `webcore:"orgId"`
	CodeId int64          `webcore:"codeId"`
	DataId core.NullInt64 `webcore:"dataId,optional"`
}

func getCode(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetCodeInput{}
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

	// Need to do a check that the user actually has access to the resource
	// that wraps the code.
	if inputs.DataId.NullInt64.Valid {
		ok, err := database.CheckValidCodeDataLink(inputs.CodeId, inputs.DataId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil || !ok {
			core.Warning("Invalid code data link: " + core.ErrorString(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code, err := webcore.GetManagedCodeFromGitea(inputs.CodeId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(code)
}

type AllCodeInput struct {
	OrgId  int32          `webcore:"orgId"`
	DataId core.NullInt64 `webcore:"dataId,optional"`
}

func allCode(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllCodeInput{}
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

	var code []*core.ManagedCode

	if inputs.DataId.NullInt64.Valid {
		code, err = database.AllManagedCodeForDataId(inputs.DataId.NullInt64.Int64, inputs.OrgId, role)
	} else {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		core.Warning("Failed to get managed code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(code)
}
