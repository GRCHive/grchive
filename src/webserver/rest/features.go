package rest

import (
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type EnableFeatureInputs struct {
	FeatureId core.FeatureId `json:"featureId"`
	OrgId     int32          `json:"orgId"`
}

func enableFeature(w http.ResponseWriter, r *http.Request) {
	inputs := EnableFeatureInputs{}
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

	err = database.EnableFeatureForOrganization(inputs.FeatureId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to start enable feature (DB): " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: should this process by async?
	switch inputs.FeatureId {
	case core.AutomationFeature:
		err = webcore.EnableAutomationFeature(inputs.OrgId)
	}

	if err != nil {
		// At this point we should be able to revert the link we did earlier safely
		// just so the user isn't stuck waiting forever. If this fails, oh well.
		database.ClearFeatureForOrganization(inputs.FeatureId, inputs.OrgId, role)

		core.Warning("Failed to enable feature: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.MarkFeatureAsFulfilled(inputs.FeatureId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to mark feature as fulfilled: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
