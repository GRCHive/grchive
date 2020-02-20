package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type AllControlSystemLinkInputs struct {
	ControlId core.NullInt64 `webcore:"controlId,optional"`
	SystemId  core.NullInt64 `webcore:"systemId,optional"`
	OrgId     int32          `webcore:"orgId"`
}

func allControlSystemLinks(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllControlSystemLinkInputs{}
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
		systems, err := database.FindSystemsLinkedToControl(inputs.ControlId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get linked systems: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonWriter.Encode(struct {
			Systems []*core.System
		}{
			Systems: systems,
		})
	} else if inputs.SystemId.NullInt64.Valid {
		controls, err := database.FindControlsLinkedToSystem(inputs.SystemId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get linked controls: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonWriter.Encode(struct {
			Controls []*core.Control
		}{
			Controls: controls,
		})
	} else {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
