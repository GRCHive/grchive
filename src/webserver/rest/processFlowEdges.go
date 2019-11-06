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

	inputOrg, err := database.FindOrganizationFromProcessFlowInputId(inputs.InputIoId, core.ServerRole)
	if err != nil {
		core.Warning("Can't get input org: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	outputOrg, err := database.FindOrganizationFromProcessFlowOutputId(inputs.OutputIoId, core.ServerRole)
	if err != nil {
		core.Warning("Can't get output org: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if inputOrg.Id != outputOrg.Id {
		core.Warning("Input and output orgs are not the same")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputOrg.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	edge, err := database.CreateNewProcessFlowEdge(&core.ProcessFlowEdge{
		InputIoId:  inputs.InputIoId,
		OutputIoId: inputs.OutputIoId,
	}, role)

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

	org, err := database.FindOrganizationFromEdgeId(inputs.EdgeId, core.ServerRole)
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

	err = database.DeleteProcessFlowEdgeFromId(inputs.EdgeId, role)
	if err != nil {
		core.Warning("Failed to delete edge: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(struct{}{})
}
