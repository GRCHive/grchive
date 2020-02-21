package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type AllControlFolderLinkInputs struct {
	ControlId core.NullInt64 `webcore:"controlId,optional"`
	FolderId  core.NullInt64 `webcore:"folderId,optional"`
	OrgId     int32          `webcore:"orgId"`
}

func allControlFolderLinks(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllControlFolderLinkInputs{}
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

	if inputs.ControlId.NullInt64.Valid {
		folders, err := database.FindFoldersLinkedToControl(inputs.ControlId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get linked folders: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonWriter.Encode(struct {
			Folders []*core.FileFolder
		}{
			Folders: folders,
		})
	} else {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
