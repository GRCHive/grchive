package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type NewNodeGLLinkInputs struct {
	NodeId    int64 `json:"nodeId"`
	AccountId int64 `json:"accountId"`
	OrgId     int32 `json:"orgId"`
}

func newNodeGLLink(w http.ResponseWriter, r *http.Request) {
	inputs := NewNodeGLLinkInputs{}
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

	err = database.NewNodeGLLink(inputs.NodeId, inputs.AccountId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to link node and GL account: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type AllNodeGLLinkInputs struct {
	NodeId    core.NullInt64 `webcore:"nodeId,optional"`
	AccountId core.NullInt64 `webcore:"accountId,optional"`
	OrgId     int32          `webcore:"orgId"`
}

func allNodeGLLink(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllNodeGLLinkInputs{}
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
		accounts, err := database.AllGLLinkedToNode(inputs.NodeId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to retrieve GL accounts linked to node: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cats, err := database.GetOrgGLCategories(inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to retrieve GL categories: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonWriter.Encode(struct {
			Accounts   []*core.GeneralLedgerAccount
			Categories []*core.GeneralLedgerCategory
		}{
			Accounts:   accounts,
			Categories: cats,
		})
	} else {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type DeleteNodeGLLinkInputs struct {
	NodeId    int64 `json:"nodeId"`
	AccountId int64 `json:"accountId"`
	OrgId     int32 `json:"orgId"`
}

func deleteNodeGLLink(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteNodeGLLinkInputs{}
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

	err = database.DeleteNodeGLLink(inputs.NodeId, inputs.AccountId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete node GL link: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
