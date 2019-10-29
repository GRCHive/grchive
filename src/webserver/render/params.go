package render

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

func BuildTemplateParams(w http.ResponseWriter, r *http.Request, needCsrf bool) map[string]interface{} {
	params := core.StructToMap(*core.LoadTemplateConfig())
	if needCsrf {
		params, _ = webcore.AddCSRFTokenToRequest(w, r, params)
	}

	_, err := webcore.FindSessionInContext(r.Context())
	params["HasSession"] = (err == nil)
	params["Host"] = r.Host
	return params
}

func BuildOrgTemplateParams(org *core.Organization) map[string]interface{} {
	params := make(map[string]interface{})
	params["OrgUrl"] = webcore.MustGetRouteUrl(
		webcore.DashboardOrgHomeRouteName,
		core.DashboardOrgOrgQueryId,
		org.OktaGroupName)
	params["OrgName"] = org.Name
	params["OrgGroupId"] = org.OktaGroupName
	return params
}

func BuildUserTemplateParams(user *core.User) map[string]interface{} {
	params := make(map[string]interface{})
	params["User"] = core.StructToMap(*user)
	return params
}

func BuildFullRiskTemplateParams(risk *core.Risk, relevantNodes []*core.ProcessFlowNode, relevantControls []*core.Control) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	rawData, err := json.Marshal(struct {
		Risk     *core.Risk
		Nodes    []*core.ProcessFlowNode
		Controls []*core.Control
	}{
		Risk:     risk,
		Nodes:    relevantNodes,
		Controls: relevantControls,
	})

	if err != nil {
		return nil, err
	}

	params["FullRiskData"] = string(rawData)
	return params, nil
}

func CreateRedirectParams(w http.ResponseWriter, r *http.Request, title string, subtitle string, redirectUrl string) map[string]interface{} {
	newMap := BuildTemplateParams(w, r, false)
	newMap["Title"] = title
	newMap["Subtitle"] = subtitle
	newMap["Redirect"] = redirectUrl
	return newMap
}
