package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
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
	OrgName           string         `webcore:"orgName"`
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
	OrgName           string         `webcore:"orgName"`
}

type GetAllControlsInputs struct {
	OrgName string `webcore:"orgName"`
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

type GetSingleControlInputs struct {
	OrgId int32 `webcore:"orgId"`
}

type LinkControlToDocCatInputs struct {
	ControlId int64 `json:"controlId"`
	OrgId     int32 `json:"orgId"`
	CatId     int64 `json:"catId"`
	IsInput   bool  `json:"isInput"`
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

	err = database.EditControl(&control, role)
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

	control := core.Control{
		Name:              inputs.Name,
		Description:       inputs.Description,
		ControlTypeId:     inputs.ControlTypeId,
		OrgId:             org.Id,
		FrequencyType:     inputs.FrequencyType,
		FrequencyInterval: inputs.FrequencyInterval,
		OwnerId:           inputs.OwnerId,
	}

	err = database.InsertNewControl(&control, role)
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
		err = database.AddControlsToNode(inputs.NodeId, []int64{control.Id}, role)
		if err != nil {
			core.Warning("Failed to add control to node relationship: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	}

	if inputs.RiskId != -1 {
		err = database.AddControlsToRisk(inputs.RiskId, []int64{control.Id}, role)
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

	apiKey, err := webcore.GetAPIKeyFromRequest(r)
	if apiKey == nil || err != nil {
		core.Warning("No API Key: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	types, err := database.GetControlTypes(core.ServerRole)
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

	if len(inputs.ControlIds) == 0 {
		return
	}

	org, err := database.FindOrganizationFromControlId(inputs.ControlIds[0], core.ServerRole)
	if err != nil {
		core.Warning("Bad organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = database.DeleteControls(
		inputs.NodeId,
		inputs.ControlIds,
		inputs.RiskIds,
		inputs.Global,
		org.Id,
		role)
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
		org, err := database.FindOrganizationFromNodeId(inputs.NodeId, core.ServerRole)
		if err != nil {
			core.Warning("Bad organization: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		role, err := webcore.GetCurrentRequestRole(r, org.Id)
		if err != nil {
			core.Warning("Bad access: " + err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = database.AddControlsToNode(inputs.NodeId, inputs.ControlIds, role)
		if err != nil {
			core.Warning("Failed to add control to node relationship: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			jsonWriter.Encode(struct{}{})
			return
		}
	}

	if inputs.RiskId != -1 {
		org, err := database.FindOrganizationFromRiskId(inputs.RiskId, core.ServerRole)
		if err != nil {
			core.Warning("Bad organization: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		role, err := webcore.GetCurrentRequestRole(r, org.Id)
		if err != nil {
			core.Warning("Bad access: " + err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = database.AddControlsToRisk(inputs.RiskId, inputs.ControlIds, role)
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

	inputs := GetAllControlsInputs{}
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

	controls, err := database.FindAllControlsForOrganization(org, role)
	if err != nil {
		core.Warning("Could not find controls: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	jsonWriter.Encode(controls)
}

func getSingleControl(w http.ResponseWriter, r *http.Request) {
	var err error
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetSingleControlInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type FullControlData struct {
		Control       *core.Control
		Nodes         []*core.ProcessFlowNode
		Risks         []*core.Risk
		InputDocCats  []*core.ControlDocumentationCategory
		OutputDocCats []*core.ControlDocumentationCategory
	}
	data := FullControlData{}
	data.Control, err = webcore.GetControlFromRequestUrl(r, core.ServerRole)
	if err != nil {
		core.Warning("No control data: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data.Nodes, err = database.FindNodesRelatedToControl(data.Control.Id, role)
	if err != nil {
		core.Warning("Failed to get nodes data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data.Risks, err = database.FindRisksRelatedToControl(data.Control.Id, role)
	if err != nil {
		core.Warning("Failed to get risks data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data.InputDocCats, err = database.GetInputDocumentationCategoriesForControl(data.Control.Id, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get input doc cats: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data.OutputDocCats, err = database.GetOutputDocumentationCategoriesForControl(data.Control.Id, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get output doc cats: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(data)
}

func linkControlToDocumentCategory(w http.ResponseWriter, r *http.Request) {
	inputs := LinkControlToDocCatInputs{}
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

	err = database.AddControlDocCatToControl(inputs.ControlId, inputs.CatId, inputs.OrgId, inputs.IsInput, role)
	if err != nil {
		core.Warning("Failed to add doc cat/control relationship: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func unlinkControlToDocumentCategory(w http.ResponseWriter, r *http.Request) {
	inputs := LinkControlToDocCatInputs{}
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

	err = database.RemoveControlDocCatFromControl(inputs.ControlId, inputs.CatId, inputs.OrgId, inputs.IsInput, role)
	if err != nil {
		core.Warning("Failed to remove doc cat/control relationship: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
