package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
	"time"
)

type GetAllProcessFlowInputs struct {
	RequestedIndex int64  `webcore:"requested"`
	OrgName        string `webcore:"organization"`
}

type DeleteProcessFlowInputs struct {
	FlowId int64 `webcore:"flowId"`
}

func getAllProcessFlows(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetAllProcessFlowInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	organization, err := database.FindOrganizationFromGroupName(inputs.OrgName)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	var flows []*core.ProcessFlow
	var index int = 0

	if inputs.RequestedIndex == -1 {
		flows, err = database.FindOrganizationProcessFlows(organization)
		if err != nil {
			core.Warning("Database error [0]: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	} else {
		flows, index, err = database.FindOrganizationProcessFlowsWithIndex(organization, inputs.RequestedIndex)
		if err != nil {
			core.Warning("Database error [1]: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	}

	jsonWriter.Encode(struct {
		Flows          []*core.ProcessFlow
		RequestedIndex int
	}{
		Flows:          flows,
		RequestedIndex: index,
	})
}

func newProcessFlow(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	// Retrieve name, description, and organization ID from the post data.
	if err := r.ParseForm(); err != nil || len(r.PostForm) == 0 {
		core.Warning("Failed to parse form data: " + core.ErrorString(err))
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	nameData := r.PostForm["name"]
	descriptionData := r.PostForm["description"]
	orgIdData := r.PostForm["organization"]

	if len(nameData) == 0 || len(descriptionData) == 0 || len(orgIdData) == 0 {
		core.Warning("Empty name or description or organization.")
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	org, err := database.FindOrganizationFromGroupName(orgIdData[0])
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	newFlow := core.ProcessFlow{
		Name:            nameData[0],
		Org:             org,
		Description:     descriptionData[0],
		CreationTime:    time.Now(),
		LastUpdatedTime: time.Now(),
	}

	err = database.InsertNewProcessFlow(&newFlow)
	if err != nil {
		if database.IsDuplicateDBEntry(err) {
			w.WriteHeader(http.StatusBadRequest)
			jsonWriter.Encode(struct {
				IsDuplicate bool
			}{
				true,
			})
		} else {
			core.Warning("Database error: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
		}
		return
	}

	// Echo the ID and Name back to the requester.
	w.WriteHeader(http.StatusOK)
	jsonWriter.Encode(struct {
		Name string
		Id   int64
	}{
		newFlow.Name,
		newFlow.Id,
	})
}

func updateProcessFlow(w http.ResponseWriter, r *http.Request) {
	// Get which process flow we want to edit, ensure user has
	// acess to it. If update successful, returns the full process
	// flow data structure (core.ProcessFlow) back to the user.
	jsonWriter := json.NewEncoder(w)

	flowId, err := webcore.GetProcessFlowIdFromRequest(r)
	if err != nil {
		core.Warning("Failed to extract flow id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	processFlow, err := database.FindProcessFlowWithId(flowId)
	if err != nil {
		core.Warning("Bad process flow id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	if err = r.ParseForm(); err != nil || len(r.PostForm) == 0 {
		core.Warning("Failed to parse form data: " + core.ErrorString(err))
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	// Only expose the name and description for editing.
	nameData := r.PostForm["name"]
	descriptionData := r.PostForm["description"]
	if len(nameData) == 0 || len(descriptionData) == 0 {
		core.Warning("Empty name or description.")
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	processFlow.Name = nameData[0]
	processFlow.Description = descriptionData[0]
	if err = database.UpdateProcessFlow(processFlow); err != nil {
		core.Warning("Failed to update flow: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonWriter.Encode(processFlow)
}

func getProcessFlowFullData(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	flowId, err := webcore.GetProcessFlowIdFromRequest(r)
	if err != nil {
		core.Warning("No flow id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	org, err := database.FindOrganizationFromProcessFlowId(flowId)
	if err != nil {
		core.Warning("Can't find organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	graph := core.ProcessFlowGraph{}
	graph.Nodes, err = database.FindAllNodesForProcessFlow(flowId)
	if err != nil {
		core.Warning("Failed to obtain nodes: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	graph.Edges, err = database.FindAllEdgesForProcessFlow(flowId)
	if err != nil {
		core.Warning("Failed to obtain edges: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	graph.Risks, err = database.FindAllRiskForOrganization(org)
	if err != nil {
		core.Warning("Failed to obtain risks: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	graph.Controls, err = database.FindAllControlsForOrganization(org)
	if err != nil {
		core.Warning("Failed to obtain controls: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	graph.NodeRisk, err = database.FindNodeRiskRelationships(flowId)
	if err != nil {
		core.Warning("Failed to obtain node-risk relationships: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	graph.NodeControl, err = database.FindNodeControlRelationships(flowId)
	if err != nil {
		core.Warning("Failed to obtain node-control relationships: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	graph.RiskControl, err = database.FindRiskControlRelationships(org.Id)
	if err != nil {
		core.Warning("Failed to obtain risk-control relationships: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(&graph)
}

func deleteProcessFlow(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := DeleteProcessFlowInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = database.DeleteProcessFlow(inputs.FlowId)
	if err != nil {
		core.Warning("Can't delete flow: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct{}{})
}
