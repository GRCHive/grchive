package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

func allDocRequestControlFolderLinks(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	req, err := webcore.FindDocumentRequestInContext(ctx)
	if err != nil {
		core.Warning("Couldn't find doc request in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	control, err := webcore.FindControlInContext(ctx)
	if err != nil {
		core.Warning("Couldn't find doc request in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	folder, err := database.FindFolderLinkedToDocRequestControl(req.Id, control.Id, req.OrgId)
	if err != nil {
		core.Warning("Failed to find folder: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(folder)
}
