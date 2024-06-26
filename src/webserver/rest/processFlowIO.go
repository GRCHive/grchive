package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"strconv"
)

type DeleteProcessFlowIOInputs struct {
	IoId    int64 `webcore:"ioId"`
	IsInput bool  `webcore:"isInput"`
}

type EditProcessFlowIOInputs struct {
	IoId    int64  `webcore:"ioId"`
	IsInput bool   `webcore:"isInput"`
	Name    string `webcore:"name"`
	Type    int32  `webcore:"type"`
}

type OrderProcessFlowIOInputs struct {
	IoId      int64 `json:"ioId"`
	IsInput   bool  `json:"isInput"`
	Direction int32 `json:"direction"`
}

func getAllProcessFlowIOTypes(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	types, err := database.GetAllProcessFlowIOTypes(core.ServerRole)
	if err != nil {
		core.Warning("Can't get IO types: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	jsonWriter.Encode(types)
}

func createNewProcessFlowIO(w http.ResponseWriter, r *http.Request) {
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
	nodeIdData := r.PostForm["nodeId"]
	isInputData := r.PostForm["isInput"]
	nameData := r.PostForm["name"]

	if len(typeIdData) == 0 || len(nodeIdData) == 0 || len(nameData) == 0 {
		core.Warning("Empty type id or node id or name.")
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

	nodeId, err := strconv.ParseInt(nodeIdData[0], 10, 64)
	if err != nil {
		core.Warning("Bad node Id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	var isInput bool = false
	if len(isInputData) != 0 {
		isInput, err = strconv.ParseBool(isInputData[0])
		if err != nil {
			core.Warning("Bad is input : " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			jsonWriter.Encode(struct{}{})
			return
		}
	}

	organization, err := database.FindOrganizationFromNodeId(nodeId, core.ServerRole)
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

	io, err := database.CreateNewProcessFlowIO(&core.ProcessFlowInputOutput{
		Id:           -1,
		Name:         nameData[0],
		ParentNodeId: nodeId,
		TypeId:       int32(typeId),
	}, isInput, role)

	if err != nil {
		core.Warning("Failed to add process flow IO: " + core.ErrorString(err))
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(io)
}

func deleteProcessFlowIO(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := DeleteProcessFlowIOInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	var organization *core.Organization
	if inputs.IsInput {
		organization, err = database.FindOrganizationFromProcessFlowInputId(inputs.IoId, core.ServerRole)
		if err != nil {
			core.Warning("Can't get input org: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		organization, err = database.FindOrganizationFromProcessFlowOutputId(inputs.IoId, core.ServerRole)
		if err != nil {
			core.Warning("Can't get output org: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	role, err := webcore.GetCurrentRequestRole(r, organization.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = database.DeleteProcessFlowIO(inputs.IoId, inputs.IsInput, role)
	if err != nil {
		core.Warning("Failed to delete IO: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(struct{}{})
}

func editProcessFlowIO(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := EditProcessFlowIOInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	var organization *core.Organization
	if inputs.IsInput {
		organization, err = database.FindOrganizationFromProcessFlowInputId(inputs.IoId, core.ServerRole)
		if err != nil {
			core.Warning("Can't get input org: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		organization, err = database.FindOrganizationFromProcessFlowOutputId(inputs.IoId, core.ServerRole)
		if err != nil {
			core.Warning("Can't get output org: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	role, err := webcore.GetCurrentRequestRole(r, organization.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	io, err := database.EditProcessFlowIO(&core.ProcessFlowInputOutput{
		Id:   inputs.IoId,
		Name: inputs.Name,
		// This doesn't need a valid value since we'll assume it'll never be updated.
		ParentNodeId: 0,
		TypeId:       inputs.Type,
	}, inputs.IsInput, role)

	if err != nil {
		core.Warning("Failed to update IO: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(io)
}

func orderProcessFlowIO(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := OrderProcessFlowIOInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	var organization *core.Organization
	if inputs.IsInput {
		organization, err = database.FindOrganizationFromProcessFlowInputId(inputs.IoId, core.ServerRole)
		if err != nil {
			core.Warning("Can't get input org: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		organization, err = database.FindOrganizationFromProcessFlowOutputId(inputs.IoId, core.ServerRole)
		if err != nil {
			core.Warning("Can't get output org: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	role, err := webcore.GetCurrentRequestRole(r, organization.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// First find the input/output in question.
	// Then we need to find the input/output that we want to swap IoOrder with.
	// Perform a swap.
	currentIo, err := database.GetProcessFlowIOFromId(inputs.IoId, inputs.IsInput, role)
	if err != nil {
		core.Warning("Failed to find current IO: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	swapIo, err := database.GetSwapProcessFlowIO(currentIo, inputs.IsInput, inputs.Direction, role)
	if err != nil {
		core.Warning("Failed to find swap IO: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.SwapIOOrder(currentIo, swapIo, inputs.IsInput, role)
	if err != nil {
		core.Warning("Failed to swap IO: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	currentIo.IoOrder, swapIo.IoOrder = swapIo.IoOrder, currentIo.IoOrder

	jsonWriter.Encode(struct {
		This  *core.ProcessFlowInputOutput
		Other *core.ProcessFlowInputOutput
	}{
		This:  currentIo,
		Other: swapIo,
	})
}
