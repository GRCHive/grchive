package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type GetResourceInputs struct {
	OrgId        int32  `webcore:"orgId"`
	ResourceType string `webcore:"resourceType"`
	ResourceId   int64  `webcore:"resourceId"`
}

func getResource(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetResourceInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	handle, err := webcore.GetResourceHandle(inputs.ResourceType, inputs.ResourceId, inputs.OrgId)
	if err != nil {
		core.Warning("Failed to get handle: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(*handle)
}
