package render

import (
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
	return params
}

func BuildOrgTemplateParams(org *core.Organization) map[string]interface{} {
	params := make(map[string]interface{})
	params["OrgUrl"] = webcore.MustGetRouteUrl(
		webcore.DashboardOrgHomeRouteName,
		core.DashboardOrgOrgQueryId,
		org.OktaGroupName)
	params["OrgName"] = org.Name
	return params
}

func BuildUserTemplateParams(user *core.User) map[string]interface{} {
	params := make(map[string]interface{})
	params["User"] = core.StructToMap(*user)
	return params
}

func CreateRedirectParams(w http.ResponseWriter, r *http.Request, title string, subtitle string, redirectUrl string) map[string]interface{} {
	newMap := BuildTemplateParams(w, r, false)
	newMap["Title"] = title
	newMap["Subtitle"] = title
	newMap["Redirect"] = redirectUrl
	return newMap
}
