package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
	"strconv"
)

type EditProcessFlowNodeInputs struct {
	NodeId      int64  `webcore:"nodeId"`
	Name        string `webcore:"name"`
	Description string `webcore:"description"`
	Type        int32  `webcore:"type"`
}

func getAllProcessFlowNodeTypes(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	types, err := database.GetAllProcessFlowNodeTypes()
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

	// Let SQL check the two IDs for validity (foreign keys) so don't bother
	// doing extra SQL queries to make sure the user input a valid type/flow.
	node, err := database.CreateNewProcessFlowNodeWithTypeId(int32(typeId), flowId)
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

	node, err := database.EditProcessFlowNode(&core.ProcessFlowNode{
		Id:            inputs.NodeId,
		Name:          inputs.Name,
		ProcessFlowId: 0,
		Description:   inputs.Description,
		NodeTypeId:    inputs.Type,
		Inputs:        nil,
		Outputs:       nil,
	})

	if err != nil {
		core.Warning("Failed to edit node: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(node)
}
