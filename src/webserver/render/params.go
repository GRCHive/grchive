package render

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type PageTemplateParameters struct {
	Organization struct {
		*core.Organization
		Url string
	} `json:"organization"`

	User struct {
		*core.User
		Auth bool
	} `json:"user"`

	Site struct {
		core.CompanyConfig
		Host string
	} `json:"site"`
}

func BuildPageTemplateParametersFull(r *http.Request) PageTemplateParameters {
	retParams := PageTemplateParameters{}
	parsedData, err := webcore.FindSessionParsedDataInContext(r.Context())

	retParams.User.Auth = (err == nil)
	if err == nil {
		retParams.User.User = parsedData.CurrentUser

		retParams.Organization.Organization = parsedData.Org
		retParams.Organization.Url = webcore.MustGetRouteUrl(
			webcore.DashboardOrgHomeRouteName,
			core.DashboardOrgOrgQueryId,
			parsedData.Org.OktaGroupName)
	}

	retParams.Site.CompanyConfig = *core.EnvConfig.Company
	retParams.Site.Host = r.Host
	return retParams
}

func BuildTemplateParams(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	params := core.StructToMap(*core.EnvConfig.Company)
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
	params["OrgId"] = org.Id
	return params
}

func BuildUserTemplateParams(user *core.User) map[string]interface{} {
	params := make(map[string]interface{})
	params["User"] = core.StructToMap(*user)
	return params
}

func CreateRedirectParams(w http.ResponseWriter, r *http.Request, title string, subtitle string, redirectUrl string) map[string]interface{} {
	newMap := BuildTemplateParams(w, r)
	newMap["Title"] = title
	newMap["Subtitle"] = subtitle
	newMap["Redirect"] = redirectUrl
	return newMap
}
