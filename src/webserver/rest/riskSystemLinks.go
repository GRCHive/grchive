package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type AllRiskSystemLinkInputs struct {
	RiskId   core.NullInt64 `webcore:"riskId,optional"`
	SystemId core.NullInt64 `webcore:"systemId,optional"`
	OrgId    int32          `webcore:"orgId"`
}

func allRiskSystemLinks(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllRiskSystemLinkInputs{}
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

	if inputs.RiskId.NullInt64.Valid {
		systems, err := database.FindSystemsLinkedToRisk(inputs.RiskId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get linked systems: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonWriter.Encode(struct {
			Systems []*core.System
		}{
			Systems: systems,
		})
	} else if inputs.SystemId.NullInt64.Valid {
		risks, err := database.FindRisksLinkedToSystem(inputs.SystemId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get linked risks: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonWriter.Encode(struct {
			Risks []*core.Risk
		}{
			Risks: risks,
		})
	} else {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
