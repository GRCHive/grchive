package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type NewRiskInputs struct {
	Name        string `webcore:"name"`
	Description string `webcore:"description"`
	NodeId      int64  `webcore:"nodeId"`
	OrgName     string `webcore:"orgName"`
}

type EditRiskInputs struct {
	Name        string `webcore:"name"`
	Description string `webcore:"description"`
	RiskId      int64  `webcore:"riskId"`
}

type DeleteRiskInputs struct {
	NodeId  int64   `webcore:"nodeId"`
	RiskIds []int64 `webcore:"riskIds"`
	Global  bool    `webcore:"global"`
}

type AddRisksToNodeInputs struct {
	NodeId  int64   `webcore:"nodeId"`
	RiskIds []int64 `webcore:"riskIds"`
}

type GetAllRisksInput struct {
	OrgName string              `webcore:"orgName"`
	Filter  core.RiskFilterData `webcore:"filter"`
}

type GetRiskInput struct {
	OrgId  int32 `webcore:"orgId"`
	RiskId int64 `webcore:"riskId"`
}

func editRisk(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := EditRiskInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	org, err := database.FindOrganizationFromRiskId(inputs.RiskId, core.ServerRole)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	risk := core.Risk{
		Id:          inputs.RiskId,
		Name:        inputs.Name,
		Description: inputs.Description,
	}
	err = database.EditRisk(&risk, role)
	if err != nil {
		core.Warning("Couldn't edit  risk: " + err.Error())
		if database.IsDuplicateDBEntry(err) {
			w.WriteHeader(http.StatusBadRequest)
			jsonWriter.Encode(database.DuplicateEntryJson)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
		}
		return
	}
	jsonWriter.Encode(risk)
}

func createNewRisk(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewRiskInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	org, err := database.FindOrganizationFromGroupName(inputs.OrgName)
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

	newRisk := core.Risk{
		Name:        inputs.Name,
		Description: inputs.Description,
		OrgId:       org.Id,
	}

	err = database.InsertNewRisk(&newRisk, role)
	if err != nil {
		core.Warning("Couldn't insert new risk: " + err.Error())
		if database.IsDuplicateDBEntry(err) {
			w.WriteHeader(http.StatusBadRequest)
			jsonWriter.Encode(database.DuplicateEntryJson)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
		}
		return
	}

	if inputs.NodeId != -1 {
		err = database.AddRisksToNode([]int64{newRisk.Id}, inputs.NodeId, role)
		if err != nil {
			core.Warning("Couldn't add risk-node relationship: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	}

	jsonWriter.Encode(newRisk)
}

func deleteRisks(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := DeleteRiskInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	if len(inputs.RiskIds) == 0 {
		return
	}

	// Assume every risk has the same organization ID.
	org, err := database.FindOrganizationFromRiskId(inputs.RiskIds[0], core.ServerRole)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = database.DeleteRisks(inputs.NodeId, inputs.RiskIds, inputs.Global, org.Id, role)
	if err != nil {
		core.Warning("Could not delete risks: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(struct{}{})
}

func addRisksToNode(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AddRisksToNodeInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	org, err := database.FindOrganizationFromNodeId(inputs.NodeId, core.ServerRole)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = database.AddRisksToNode(inputs.RiskIds, inputs.NodeId, role)
	if err != nil {
		core.Warning("Couldn't add risks: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(struct{}{})
}

func getAllRisks(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetAllRisksInput{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromGroupName(inputs.OrgName)
	if err != nil {
		core.Warning("No organization data: " + err.Error())
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

	risks, err := database.FindAllRiskForOrganization(org, inputs.Filter, role)
	if err != nil {
		core.Warning("Could not find risks: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(risks)
}

func getSingleRisk(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	type FullRiskData struct {
		Risk     *core.Risk
		Flows    []*core.ProcessFlow
		Controls []*core.Control
	}

	inputs := GetRiskInput{}
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

	data := FullRiskData{}
	data.Risk, err = database.FindRisk(inputs.RiskId, role)
	if err != nil {
		core.Warning("No risk data: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	data.Flows, err = database.FindFlowsRelatedToRisk(data.Risk.Id, role)
	if err != nil {
		core.Warning("Failed to get nodes data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	data.Controls, err = database.FindControlsRelatedToRisk(data.Risk.Id, role)
	if err != nil {
		core.Warning("Failed to get controls data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(data)
}
