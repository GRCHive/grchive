package render

import (
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type PageTemplateParameters struct {
	Organization struct {
		*core.Organization
		Url string
	} `json:"organization"`

	User struct {
		*core.User
		Auth     bool
		Verified bool
	} `json:"user"`

	Site struct {
		core.CompanyConfig
		Host string
	} `json:"site"`

	Auth struct {
		OktaServer      string
		OktaClientId    string
		OktaRedirectUri string
		OktaScope       string
	} `json:"auth"`
}

func BuildPageTemplateParametersFull(r *http.Request) PageTemplateParameters {
	retParams := PageTemplateParameters{}
	parsedData, err := webcore.FindSessionParsedDataInContext(r.Context())

	retParams.User.Auth = (err == nil)
	retParams.User.Verified = false

	if err == nil {
		retParams.User.User = parsedData.CurrentUser
		retParams.User.Verified = parsedData.VerifiedEmail
	}

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err == nil {
		retParams.Organization.Organization = org
		retParams.Organization.Url = webcore.MustGetRouteUrl(
			webcore.DashboardOrgHomeRouteName,
			core.DashboardOrgOrgQueryId,
			org.OktaGroupName)
	}

	retParams.Site.CompanyConfig = *core.EnvConfig.Company
	retParams.Site.Host = r.Host

	retParams.Auth.OktaServer = core.EnvConfig.Okta.BaseUrl
	retParams.Auth.OktaClientId = core.EnvConfig.Login.ClientId
	retParams.Auth.OktaRedirectUri = webcore.MustGetRouteUrlAbsolute(webcore.SamlCallbackRouteName)
	retParams.Auth.OktaScope = core.EnvConfig.Login.Scope
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
