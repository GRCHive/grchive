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

	newRisk := core.Risk{
		Name:            inputs.Name,
		Description:     inputs.Description,
		RelevantNodeIds: []int64{inputs.NodeId},
	}

	err = database.InsertNewRisk(&newRisk)
	if err != nil {
		core.Warning("Couldn't insert new risk: " + err.Error())
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
