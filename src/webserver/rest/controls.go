package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type NewControlInputs struct {
	Name              string         `webcore:"name"`
	Description       string         `webcore:"description"`
	ControlTypeId     int32          `webcore:"controlType"`
	FrequencyType     int32          `webcore:"frequencyType"`
	FrequencyInterval int32          `webcore:"frequencyInterval"`
	OwnerId           core.NullInt64 `webcore:"ownerId,optional"`
	NodeId            int64          `webcore:"nodeId"`
	RiskId            int64          `webcore:"riskId"`
}

type EditControlInputs struct {
	Name              string         `webcore:"name"`
	Description       string         `webcore:"description"`
	ControlTypeId     int32          `webcore:"controlType"`
	FrequencyType     int32          `webcore:"frequencyType"`
	FrequencyInterval int32          `webcore:"frequencyInterval"`
	OwnerId           core.NullInt64 `webcore:"ownerId,optional"`
	NodeId            int64          `webcore:"nodeId"`
	RiskId            int64          `webcore:"riskId"`
	ControlId         int64          `webcore:"controlId"`
}

type DeleteControlInputs struct {
	NodeId     int64   `webcore:"nodeId"`
	RiskIds    []int64 `webcore:"riskIds"`
	ControlIds []int64 `webcore:"controlIds"`
	Global     bool    `webcore:"global"`
}

type AddControlInputs struct {
	NodeId     int64   `webcore:"nodeId"`
	RiskId     int64   `webcore:"riskId"`
	ControlIds []int64 `webcore:"controlIds"`
}

func editControl(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := EditControlInputs{}
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

	control := core.Control{
		Id:                inputs.ControlId,
		Name:              inputs.Name,
		Description:       inputs.Description,
		ControlTypeId:     inputs.ControlTypeId,
		OrgId:             org.Id,
		FrequencyType:     inputs.FrequencyType,
		FrequencyInterval: inputs.FrequencyInterval,
		OwnerId:           inputs.OwnerId,
	}

	err = database.EditControl(&control)
	if err != nil {
		core.Warning("Failed to edit control: " + err.Error())
		if database.IsDuplicateDBEntry(err) {
			w.WriteHeader(http.StatusBadRequest)
			jsonWriter.Encode(database.DuplicateEntryJson)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
		}
		return
	}

	jsonWriter.Encode(control)
}

func newControl(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewControlInputs{}
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

	control := core.Control{
		Name:              inputs.Name,
		Description:       inputs.Description,
		ControlTypeId:     inputs.ControlTypeId,
		OrgId:             org.Id,
		FrequencyType:     inputs.FrequencyType,
		FrequencyInterval: inputs.FrequencyInterval,
		OwnerId:           inputs.OwnerId,
	}

	err = database.InsertNewControl(&control)
	if err != nil {
		core.Warning("Failed to insert new control: " + err.Error())
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
		err = database.AddControlsToNode(inputs.NodeId, []int64{control.Id})
		if err != nil {
			core.Warning("Failed to add control to node relationship: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	}

	if inputs.RiskId != -1 {
		err = database.AddControlsToRisk(inputs.RiskId, []int64{control.Id})
		if err != nil {
			core.Warning("Failed to add control to risk relationship: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	}

	jsonWriter.Encode(control)
}

func getControlTypes(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	types, err := database.GetControlTypes()
	if err != nil {
		core.Warning("Failed to get control types: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(types)
}

func deleteControls(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := DeleteControlInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = database.DeleteControls(
		inputs.NodeId,
		inputs.ControlIds,
		inputs.RiskIds,
		inputs.Global)
	if err != nil {
		core.Warning("Failed to delete controls: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonWriter.Encode(struct{}{})
}

func addControls(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AddControlInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if inputs.NodeId != -1 {
		err = database.AddControlsToNode(inputs.NodeId, inputs.ControlIds)
		if err != nil {
			core.Warning("Failed to add control to node relationship: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	}

	if inputs.RiskId != -1 {
		err = database.AddControlsToRisk(inputs.RiskId, inputs.ControlIds)
		if err != nil {
			core.Warning("Failed to add control to risk relationship: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	}

	jsonWriter.Encode(struct{}{})
}

func getAllControls(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	userParsedData, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("No user session data: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	controls, err := database.FindAllControlsForOrganization(userParsedData.Org)
	if err != nil {
		core.Warning("Could not find controls: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(controls)
}
