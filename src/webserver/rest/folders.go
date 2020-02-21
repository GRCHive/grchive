package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type NewFolderInputs struct {
	Name      string `json:"name"`
	OrgId     int32  `json:"orgId"`
	ControlId int64  `json:"controlId"`
}

func newFolder(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewFolderInputs{}
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

	folder := core.FileFolder{
		OrgId: inputs.OrgId,
		Name:  inputs.Name,
	}

	tx := database.CreateTx()

	err = database.NewFolderWithTx(&folder, role, tx)
	if err != nil {
		tx.Rollback()
		core.Warning("Failed to create new folder: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.AddControlFolderLinkWithTx(inputs.ControlId, folder.Id, inputs.OrgId, role, tx)
	if err != nil {
		tx.Rollback()
		core.Warning("Failed to add folder control link: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		core.Warning("Failed to commit: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(folder)
}

type UpdateFolderInputs struct {
	Name     string `json:"name"`
	OrgId    int32  `json:"orgId"`
	FolderId int64  `json:"folderId"`
}

func updateFolder(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UpdateFolderInputs{}
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

	folder := core.FileFolder{
		Id:    inputs.FolderId,
		OrgId: inputs.OrgId,
		Name:  inputs.Name,
	}

	err = database.UpdateFolder(&folder, role)
	if err != nil {
		core.Warning("Failed to update folder: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(folder)
}

type DeleteFolderInputs struct {
	OrgId    int32 `json:"orgId"`
	FolderId int64 `json:"folderId"`
}

func deleteFolder(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteFolderInputs{}
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

	err = database.DeleteFolder(inputs.FolderId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete folder: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
