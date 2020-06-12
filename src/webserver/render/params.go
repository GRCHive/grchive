package render

import (
	"github.com/gorilla/mux"
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

	Resource struct {
		Id string
	} `json:"resource"`

	// These are available for server-side rendering
	ServerSide struct {
		OrgName  string
		UserName string
	} `json:"-"`
}

func BuildPageTemplateParametersFull(r *http.Request, resourceQuery string) PageTemplateParameters {
	retParams := PageTemplateParameters{}
	parsedData, err := webcore.FindSessionParsedDataInContext(r.Context())

	retParams.User.Auth = (err == nil)
	retParams.User.Verified = false

	if err == nil {
		retParams.User.User = parsedData.CurrentUser
		retParams.User.Verified = parsedData.VerifiedEmail
		retParams.ServerSide.UserName = retParams.User.User.FullName()
	}

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err == nil {
		retParams.Organization.Organization = org
		retParams.Organization.Url = webcore.MustGetRouteUrl(
			webcore.DashboardOrgHomeRouteName,
			core.DashboardOrgOrgQueryId,
			org.OktaGroupName)
		retParams.ServerSide.OrgName = org.Name
	}

	retParams.Site.CompanyConfig = *core.EnvConfig.Company
	retParams.Site.Host = r.Host

	retParams.Auth.OktaServer = core.EnvConfig.Okta.BaseUrl
	retParams.Auth.OktaClientId = core.EnvConfig.Login.ClientId
	retParams.Auth.OktaRedirectUri = webcore.MustGetRouteUrlAbsolute(webcore.SamlCallbackRouteName)
	retParams.Auth.OktaScope = core.EnvConfig.Login.Scope

	if resourceQuery != "" {
		urlRouteVars := mux.Vars(r)
		resource, ok := urlRouteVars[resourceQuery]
		if ok {
			retParams.Resource.Id = resource
		}
	}
	return retParams
}

func BuildTemplateParams(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	params := core.StructToMap(*core.EnvConfig.Company)
	_, err := webcore.FindSessionInContext(r.Context())
	params["HasSession"] = (err == nil)
	params["Host"] = r.Host
	return params
}

func CreateRedirectParams(w http.ResponseWriter, r *http.Request, title string, subtitle string, redirectUrl string) map[string]interface{} {
	newMap := BuildTemplateParams(w, r)
	newMap["Title"] = title
	newMap["Subtitle"] = subtitle
	newMap["Redirect"] = redirectUrl
	return newMap
}

func CoreParams() map[string]interface{} {
	params := make(map[string]interface{})
	params["UseAnalytics"] = core.EnvConfig.UseAnalytics
	return params
}
