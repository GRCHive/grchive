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

type DeleteRiskInputs struct {
	NodeId  int64   `webcore:"nodeId"`
	RiskIds []int64 `webcore:"riskIds"`
	Global  bool    `webcore:"global"`
}

type AddRisksToNodeInputs struct {
	NodeId  int64   `webcore:"nodeId"`
	RiskIds []int64 `webcore:"riskIds"`
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

	err = database.AddRisksToNode([]int64{newRisk.Id}, inputs.NodeId)
	if err != nil {
		core.Warning("Couldn't add risk-node relationship: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
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
