package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type NewNodeSystemLinkInputs struct {
	NodeId   int64 `json:"nodeId"`
	SystemId int64 `json:"systemId"`
	OrgId    int32 `json:"orgId"`
}

func newNodeSystemLink(w http.ResponseWriter, r *http.Request) {
	inputs := NewNodeSystemLinkInputs{}
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

	err = database.NewNodeSystemLink(inputs.NodeId, inputs.SystemId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to link node and system: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type AllNodeSystemLinkInputs struct {
	NodeId   core.NullInt64 `webcore:"nodeId,optional"`
	SystemId core.NullInt64 `webcore:"systemId,optional"`
	OrgId    int32          `webcore:"orgId"`
}

func allNodeSystemLink(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllNodeSystemLinkInputs{}
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

	if inputs.NodeId.NullInt64.Valid {
		systems, err := database.AllSystemsLinkedToNode(inputs.NodeId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to retrieve systems linked to node: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonWriter.Encode(systems)
	} else {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type DeleteNodeSystemLinkInputs struct {
	NodeId   int64 `json:"nodeId"`
	SystemId int64 `json:"systemId"`
	OrgId    int32 `json:"orgId"`
}

func deleteNodeSystemLink(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteNodeSystemLinkInputs{}
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

	err = database.DeleteNodeSystemLink(inputs.NodeId, inputs.SystemId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete node system link: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
