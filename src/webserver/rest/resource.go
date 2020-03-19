package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"strconv"
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

	name, err := database.GetResourceName(inputs.ResourceType, inputs.ResourceId)
	if err != nil {
		core.Warning("Failed to get resource handle: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("Failed to find org: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resourceIdStr := strconv.FormatInt(inputs.ResourceId, 10)

	var url core.NullString
	switch inputs.ResourceType {
	case core.ResourceControl:
		url = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
			webcore.SingleControlRouteName,
			core.DashboardOrgOrgQueryId, org.OktaGroupName,
			core.DashboardOrgControlQueryId, resourceIdStr,
		))
	case core.ResourceDocRequest:
		url = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
			webcore.SingleDocRequestRouteName,
			core.DashboardOrgOrgQueryId, org.OktaGroupName,
			core.DashboardOrgDocRequestQueryId, resourceIdStr,
		))
	case core.ResourceSqlQueryRequest:
		url = core.CreateNullString(webcore.MustGetRouteUrlAbsolute(
			webcore.SingleSqlRequestRouteName,
			core.DashboardOrgOrgQueryId, org.OktaGroupName,
			core.DashboardOrgSqlRequestQueryId, resourceIdStr,
		))
	}

	jsonWriter.Encode(core.ResourceHandle{
		DisplayText: name,
		ResourceUri: url,
	})
}
