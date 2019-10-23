package webcore

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"net/url"
)

// Central location for determining routes to things.
// Note that the main application when registering routes needs to use core.url
// but everyone else should use webcore.router after the router is finished being created.

type RouteName string

const (
	GettingStartedRouteName              RouteName = "GettingStarted"
	ContactUsRouteName                             = "ContactUs"
	LandingPageRouteName                           = "LandingPage"
	LoginRouteName                                 = "Login"
	LogoutRouteName                                = "Logout"
	LearnMoreRouteName                             = "LearnMore"
	GettingStartedPostRouteName                    = "GettingStartedPost"
	LoginPostRouteName                             = "LoginPost"
	SamlCallbackRouteName                          = "SamlCallback"
	DashboardHomeRouteName                         = "DashboardHome"
	DashboardOrgHomeRouteName                      = "DashboardOrgHome"
	DashboardUserHomeRouteName                     = "DashboardUserHome"
	UserProfileEditRouteName                       = "UserProfilePost"
	DashboardProcessFlowsRouteName                 = "ProcessFlows"
	NewProcessFlowRouteName                        = "NewProcessFlow"
	GetAllProcessFlowRouteName                     = "GetAllProcessFlow"
	UpdateProcessFlowRouteName                     = "UpdateProcessFlow"
	GetAllProcessFlowNodeTypesRouteName            = "GetAllProcessFlowNodeTypes"
	GetProcessFlowFullDataRouteName                = "GetProcessFlowFullData"
	NewProcessFlowNodeRouteName                    = "NewProcessFlowNode"
	GetAllProcessFlowIOTypesRouteName              = "GetAllProcessFlowIOTypes"
	CreateNewProcessFlowIOTypesRouteName           = "CreateNewProcessFlowIOTypes"
	DeleteProcessFlowIORouteName                   = "DeleteProcessFlowIO"
	EditProcessFlowIORouteName                     = "EditProcessFlowIO"
	EditProcessFlowNodeRouteName                   = "EditProcessFlowNode"
	CreateNewProcessFlowEdgeRouteName              = "CreateNewProcessFlowEdge"
	DeleteProcessFlowEdgeRouteName                 = "DeleteProcessFlowEdge"
	DeleteProcessFlowNodeRouteName                 = "DeleteProcessFlowNode"
	NewRiskRouteName                               = "NewRisk"
	DeleteRiskRouteName                            = "DeleteRisk"
	AddRiskToNodeRouteName                         = "AddRiskToNodeRisk"
	GetAllOrgUsersRouteName                        = "GetAllOrgUsers"
	NewControlRouteName                            = "NewControl"
	ControlTypesRouteName                          = "ControlTypes"
)

var globalRouter *mux.Router

func RegisterRouter(r *mux.Router) {
	globalRouter = r
}

func MustGetRouteUrl(r RouteName, params ...string) string {
	route := globalRouter.Get(string(r))
	if route == nil {
		core.Warning("No Route: " + string(r))
		return "/404"
	}

	url, err := route.URL(params...)
	if err != nil {
		core.Warning("Bad Route: " + string(r) + " :: " + err.Error())
		return "/404"
	}
	return url.String()
}

func MustGetRouteUrlAbsolute(r RouteName, params ...string) string {
	return core.LoadEnvConfig().SelfUri + MustGetRouteUrl(r, params...)
}

func CreateOktaLoginUrl(idp string, state string, nonce string) string {
	envConfig := core.LoadEnvConfig()

	return fmt.Sprintf("%s%s%s?idp=%s&client_id=%s&response_type=%s&response_mode=%s&scope=%s&redirect_uri=%s&state=%s&=nonce=%s",
		envConfig.Okta.BaseUrl,
		envConfig.Login.AuthServerEndpoint,
		envConfig.Login.AuthEndpoint,
		idp,
		envConfig.Login.ClientId,
		envConfig.Login.ResponseType,
		envConfig.Login.ResponseMode,
		url.QueryEscape(envConfig.Login.Scope),
		url.QueryEscape(MustGetRouteUrlAbsolute(SamlCallbackRouteName)),
		state,
		url.QueryEscape(nonce))
}

var OktaTokenUrl string = fmt.Sprintf("%s%s%s",
	core.LoadEnvConfig().Okta.BaseUrl,
	core.LoadEnvConfig().Login.AuthServerEndpoint,
	core.LoadEnvConfig().Login.TokenEndpoint)

var OktaKeyUrl string = fmt.Sprintf("%s%s%s?client_id=%s",
	core.LoadEnvConfig().Okta.BaseUrl,
	core.LoadEnvConfig().Login.AuthServerEndpoint,
	core.LoadEnvConfig().Login.KeyEndpoint,
	core.LoadEnvConfig().Login.ClientId)

func CreateOktaLogoutUrl(idToken string) string {
	envConfig := core.LoadEnvConfig()

	return fmt.Sprintf("%s%s%s?id_token_hint=%s&post_logout_redirect_uri=%s",
		envConfig.Okta.BaseUrl,
		envConfig.Login.AuthServerEndpoint,
		envConfig.Login.LogoutEndpoint,
		idToken,
		url.QueryEscape(MustGetRouteUrlAbsolute(LandingPageRouteName)))
}

func CreateOktaUserUpdateUrl(userId string) string {
	envConfig := core.LoadEnvConfig()

	return fmt.Sprintf("%s%s%s/%s",
		envConfig.Okta.BaseUrl,
		envConfig.Okta.ApiEndpoint,
		envConfig.Okta.UsersEndpoint,
		userId)
}
