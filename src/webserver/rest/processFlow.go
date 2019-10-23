package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
	"strconv"
	"time"
)

func getAllProcessFlows(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	queryVals := r.URL.Query()

	organizationName, ok := queryVals["organization"]
	if !ok || len(organizationName) == 0 {
		core.Warning("Failed to get process flows (no organization)")
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	organization, err := database.FindOrganizationFromGroupName(organizationName[0])
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	// Ensure that the user has access to this organization.
	userParsedData, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("No user session data: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	if userParsedData.Org.OktaGroupId != organization.OktaGroupId {
		core.Warning("Unauthorized access")
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	requestedId, ok := queryVals["requested"]

	var flows []*core.ProcessFlow
	var index int = 0

	if !ok || len(requestedId) == 0 {
		flows, err = database.FindOrganizationProcessFlows(organization)
		if err != nil {
			core.Warning("Database error [0]: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	} else {
		intRequestedId, err := strconv.ParseInt(requestedId[0], 10, 64)
		if err != nil {
			core.Warning("Invalid requested id: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			jsonWriter.Encode(struct{}{})
			return
		}

		flows, index, err = database.FindOrganizationProcessFlowsWithIndex(organization, int64(intRequestedId))
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

	// Ensure that the user has access to this organization.
	userParsedData, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("No user session data: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	if userParsedData.Org.OktaGroupId != org.OktaGroupId {
		core.Warning("Unauthorized access")
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

	userData, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("No user session context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	if userData.Org.OktaGroupId != processFlow.Org.OktaGroupId {
		core.Warning("Permission denied.")
		w.WriteHeader(http.StatusForbidden)
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

	org, err := webcore.FindOrganizationInContext(r.Context())
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

	jsonWriter.Encode(&graph)
}
