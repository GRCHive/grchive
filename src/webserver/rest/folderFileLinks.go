package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type AllFolderFileLinkInputs struct {
	FolderId core.NullInt64 `webcore:"folderId,optional"`
	OrgId    int32          `webcore:"orgId"`
}

func allFolderFileLinks(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllFolderFileLinkInputs{}
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

	if inputs.FolderId.NullInt64.Valid {
		files, err := database.FindFilesLinkedToFolder(inputs.FolderId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get linked files: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonWriter.Encode(struct {
			Files []*core.ControlDocumentationFile
		}{
			Files: files,
		})
	} else {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type NewFolderFileLinkInputs struct {
	FolderId int64   `json:"folderId"`
	FileIds  []int64 `json:"fileIds"`
	OrgId    int32   `json:"orgId"`
}

func newFolderFileLinks(w http.ResponseWriter, r *http.Request) {
	inputs := NewFolderFileLinkInputs{}
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

	for _, id := range inputs.FileIds {
		err = database.AddFileToFolderWithTx(id, inputs.FolderId, inputs.OrgId, role, tx)
		if err != nil {
			tx.Rollback()
			core.Warning("Failed to add file to folder: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		core.Warning("Failed to commit: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type DeleteFolderFileLinkInputs struct {
	FolderId int64 `json:"folderId"`
	FileId   int64 `json:"fileId"`
	OrgId    int32 `json:"orgId"`
}

func deleteFolderFileLink(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteFolderFileLinkInputs{}
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

	err = database.DeleteFileFromFolder(inputs.FileId, inputs.FolderId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete file from folder: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
