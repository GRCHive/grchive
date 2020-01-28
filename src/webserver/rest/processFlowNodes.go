package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"strconv"
)

type EditProcessFlowNodeInputs struct {
	NodeId      int64  `webcore:"nodeId"`
	Name        string `webcore:"name"`
	Description string `webcore:"description"`
	Type        int32  `webcore:"type"`
}

type DeleteProcessFlowNodeInputs struct {
	NodeId int64 `webcore:"nodeId"`
}

func getAllProcessFlowNodeTypes(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	apiKey, err := webcore.GetAPIKeyFromRequest(r)
	if apiKey == nil || err != nil {
		core.Warning("No API Key: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	types, err := database.GetAllProcessFlowNodeTypes(core.ServerRole)
	if err != nil {
		core.Warning("Can't get types: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	jsonWriter.Encode(types)
}

func newProcessFlowNode(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	// Retrieve name, description, and organization ID from the post data.
	if err := r.ParseForm(); err != nil || len(r.PostForm) == 0 {
		core.Warning("Failed to parse form data: " + core.ErrorString(err))
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	typeIdData := r.PostForm["typeId"]
	flowIdData := r.PostForm["flowId"]

	if len(typeIdData) == 0 || len(flowIdData) == 0 {
		core.Warning("Empty type id or flow id.")
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	typeId, err := strconv.ParseInt(typeIdData[0], 10, 32)
	if err != nil {
		core.Warning("Bad type Id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	flowId, err := strconv.ParseInt(flowIdData[0], 10, 64)
	if err != nil {
		core.Warning("Bad flow Id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	org, err := database.FindOrganizationFromProcessFlowId(flowId, core.ServerRole)
	if err != nil {
		core.Warning("Can't find organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Let SQL check the two IDs for validity (foreign keys) so don't bother
	// doing extra SQL queries to make sure the user input a valid type/flow.
	node, err := database.CreateNewProcessFlowNodeWithTypeId(int32(typeId), flowId, role)
	if err != nil {
		core.Warning("Failed to create new process flow node: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(struct {
		Node *core.ProcessFlowNode
	}{
		Node: node,
	})
}

func editProcessFlowNode(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := EditProcessFlowNodeInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	org, err := database.FindOrganizationFromNodeId(inputs.NodeId, core.ServerRole)
	if err != nil {
		core.Warning("Can't find organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	node, err := database.EditProcessFlowNode(&core.ProcessFlowNode{
		Id:            inputs.NodeId,
		Name:          inputs.Name,
		ProcessFlowId: 0,
		Description:   inputs.Description,
		NodeTypeId:    inputs.Type,
		Inputs:        nil,
		Outputs:       nil,
	}, role)

	if err != nil {
		core.Warning("Failed to edit node: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(node)
}

func deleteProcessFlowNode(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := DeleteProcessFlowNodeInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	org, err := database.FindOrganizationFromNodeId(inputs.NodeId, core.ServerRole)
	if err != nil {
		core.Warning("Can't find organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = database.DeleteProcessFlowNodeFromId(inputs.NodeId, role)
	if err != nil {
		core.Warning("Failed to delete node: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(struct{}{})
}
