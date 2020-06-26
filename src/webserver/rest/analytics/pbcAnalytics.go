package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type OverallPbcAnalyticsInputs struct {
	Filter core.DocRequestFilterData `webcore:"filter"`
}

func getOverallPbcAnalytics(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := OverallPbcAnalyticsInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := database.GetOverallPbcAnalytics(org.Id, inputs.Filter)
	if err != nil {
		core.Warning("Failed to retrieve analytics: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(data)
}

func getCategoryRequesterPbcAnalytics(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := OverallPbcAnalyticsInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := database.GetRequesterCategoryPbcAnalytics(org.Id, inputs.Filter)
	if err != nil {
		core.Warning("Failed to retrieve analytics: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(data)
}

func getCategoryAssigneePbcAnalytics(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := OverallPbcAnalyticsInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := database.GetAssigneeCategoryPbcAnalytics(org.Id, inputs.Filter)
	if err != nil {
		core.Warning("Failed to retrieve analytics: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(data)
}

func getCategoryDocCatPbcAnalytics(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := OverallPbcAnalyticsInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := database.GetDocCatCategoryPbcAnalytics(org.Id, inputs.Filter)
	if err != nil {
		core.Warning("Failed to retrieve analytics: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(data)
}

func getCategoryProcessFlowPbcAnalytics(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := OverallPbcAnalyticsInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := database.GetProcessFlowCategoryPbcAnalytics(org.Id, inputs.Filter)
	if err != nil {
		core.Warning("Failed to retrieve analytics: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(data)
}

func getCategoryControlPbcAnalytics(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := OverallPbcAnalyticsInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := database.GetControlCategoryPbcAnalytics(org.Id, inputs.Filter)
	if err != nil {
		core.Warning("Failed to retrieve analytics: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(data)
}

func getCategoryRiskPbcAnalytics(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := OverallPbcAnalyticsInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := database.GetRiskCategoryPbcAnalytics(org.Id, inputs.Filter)
	if err != nil {
		core.Warning("Failed to retrieve analytics: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(data)
}

func getCategoryGLPbcAnalytics(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := OverallPbcAnalyticsInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := database.GetGLCategoryPbcAnalytics(org.Id, inputs.Filter)
	if err != nil {
		core.Warning("Failed to retrieve analytics: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(data)
}

func getCategorySystemPbcAnalytics(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find org in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := OverallPbcAnalyticsInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := database.GetSystemCategoryPbcAnalytics(org.Id, inputs.Filter)
	if err != nil {
		core.Warning("Failed to retrieve analytics: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonWriter.Encode(data)
}
