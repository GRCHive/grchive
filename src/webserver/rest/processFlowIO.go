package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"net/http"
	"strconv"
)

func getAllProcessFlowIOTypes(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	types, err := database.GetAllProcessFlowIOTypes()
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

	io, err := database.CreateNewProcessFlowIO(&core.ProcessFlowInputOutput{
		Id:           -1,
		Name:         nameData[0],
		ParentNodeId: nodeId,
		TypeId:       int32(typeId),
	}, isInput)

	if err != nil {
		core.Warning("Failed to add process flow IO: " + core.ErrorString(err))
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(io)
}
