package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type NewProcessFlowEdgeInputs struct {
	InputIoId  int64 `webcore:"inputIoId"`
	OutputIoId int64 `webcore:"outputIoId"`
}

type DeleteProcessFlowEdgeInputs struct {
	EdgeId int64 `webcore:"edgeId"`
}

func createNewProcessFlowEdge(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewProcessFlowEdgeInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	edge, err := database.CreateNewProcessFlowEdge(&core.ProcessFlowEdge{
		InputIoId:  inputs.InputIoId,
		OutputIoId: inputs.OutputIoId,
	})

	if err != nil {
		core.Warning("Can't add edge: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(edge)
}

func deleteProcessFlowEdge(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := DeleteProcessFlowEdgeInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	err = database.DeleteEdgeFromId(inputs.EdgeId)
	if err != nil {
		core.Warning("Failed to delete edge: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(struct{}{})
}
