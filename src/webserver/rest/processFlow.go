package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
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

type UpdateProcessFlowInputs struct {
	Name        string `webcore:"name"`
	Description string `webcore:"description"`
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

	role, err := webcore.GetCurrentRequestRole(r, organization.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var flows []*core.ProcessFlow
	var index int = 0

	if inputs.RequestedIndex == -1 {
		flows, err = database.FindOrganizationProcessFlows(organization, role)
		if err != nil {
			core.Warning("Database error [0]: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	} else {
		flows, index, err = database.FindOrganizationProcessFlowsWithIndex(organization, inputs.RequestedIndex, role)
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

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newFlow := core.ProcessFlow{
		Name:            nameData[0],
		Org:             org,
		Description:     descriptionData[0],
		CreationTime:    time.Now(),
		LastUpdatedTime: time.Now(),
	}

	err = database.InsertNewProcessFlow(&newFlow, role)
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
	jsonWriter.Encode(newFlow)
}

func updateProcessFlow(w http.ResponseWriter, r *http.Request) {
	// Get which process flow we want to edit, ensure user has
	// acess to it. If update successful, returns the full process
	// flow data structure (core.ProcessFlow) back to the user.
	jsonWriter := json.NewEncoder(w)

	inputs := UpdateProcessFlowInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	flowId, err := webcore.GetProcessFlowIdFromRequest(r)
	if err != nil {
		core.Warning("Failed to extract flow id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	processFlow, err := database.FindProcessFlowWithId(flowId, core.ServerRole)
	if err != nil {
		core.Warning("Bad process flow id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, processFlow.Org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	processFlow.Name = inputs.Name
	processFlow.Description = inputs.Description
	if err = database.UpdateProcessFlow(processFlow, role); err != nil {
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

	flow, err := database.FindProcessFlowWithId(flowId, role)
	if err != nil {
		core.Warning("Failed to obtain flow data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	graph := core.ProcessFlowGraph{}
	graph.Nodes, err = database.FindAllNodesForProcessFlow(flowId, role)
	if err != nil {
		core.Warning("Failed to obtain nodes: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	graph.Edges, err = database.FindAllEdgesForProcessFlow(flowId, role)
	if err != nil {
		core.Warning("Failed to obtain edges: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	graph.Risks, err = database.FindAllRiskForOrganization(org, core.NullRiskFilterData, role)
	if err != nil {
		core.Warning("Failed to obtain risks: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	graph.Controls, err = database.FindAllControlsForOrganization(org, role)
	if err != nil {
		core.Warning("Failed to obtain controls: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	graph.NodeRisk, err = database.FindNodeRiskRelationships(flowId, role)
	if err != nil {
		core.Warning("Failed to obtain node-risk relationships: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	graph.NodeControl, err = database.FindNodeControlRelationships(flowId, role)
	if err != nil {
		core.Warning("Failed to obtain node-control relationships: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	graph.RiskControl, err = database.FindRiskControlRelationships(org.Id, role)
	if err != nil {
		core.Warning("Failed to obtain risk-control relationships: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(struct {
		Basic *core.ProcessFlow
		Graph *core.ProcessFlowGraph
	}{
		Basic: flow,
		Graph: &graph,
	})
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

	org, err := database.FindOrganizationFromProcessFlowId(inputs.FlowId, core.ServerRole)
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

	err = database.DeleteProcessFlow(inputs.FlowId, role)
	if err != nil {
		core.Warning("Can't delete flow: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct{}{})
}
