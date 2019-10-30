package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type NewRiskInputs struct {
	Name        string `webcore:"name"`
	Description string `webcore:"description"`
	NodeId      int64  `webcore:"nodeId"`
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

	risk := core.Risk{
		Id:          inputs.RiskId,
		Name:        inputs.Name,
		Description: inputs.Description,
	}
	err = database.EditRisk(&risk)
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

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	newRisk := core.Risk{
		Name:        inputs.Name,
		Description: inputs.Description,
		Org:         org,
	}

	err = database.InsertNewRisk(&newRisk)
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
		err = database.AddRisksToNode([]int64{newRisk.Id}, inputs.NodeId)
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

	err = database.DeleteRisks(inputs.NodeId, inputs.RiskIds, inputs.Global)
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

	err = database.AddRisksToNode(inputs.RiskIds, inputs.NodeId)
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

	userParsedData, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("No user session data: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	risks, err := database.FindAllRiskForOrganization(userParsedData.Org)
	if err != nil {
		core.Warning("Could not find risks: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(risks)
}

func getSingleRisk(w http.ResponseWriter, r *http.Request) {
	var err error
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	type FullRiskData struct {
		Risk     *core.Risk
		Nodes    []*core.ProcessFlowNode
		Controls []*core.Control
	}
	data := FullRiskData{}
	data.Risk, err = webcore.GetRiskFromRequestUrl(r)
	if err != nil {
		core.Warning("No risk data: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	data.Nodes, err = database.FindNodesRelatedToRisk(data.Risk.Id)
	if err != nil {
		core.Warning("Failed to get nodes data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	data.Controls, err = database.FindControlsRelatedToRisk(data.Risk.Id)
	if err != nil {
		core.Warning("Failed to get controls data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(data)
}
