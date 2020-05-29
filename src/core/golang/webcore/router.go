package webcore

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/grchive/grchive/core"
	"net/url"
)

// Central location for determining routes to things.
// Note that the main application when registering routes needs to use core.url
// but everyone else should use webcore.router after the router is finished being created.

type RouteName string

const (
	LandingPageRouteName           = "LandingPage"
	SamlCallbackRouteName          = "SamlCallback"
	EmailVerifyRouteName           = "VerifyEmail"
	AcceptInviteRouteName          = "AcceptInvite"
	RegisterRouteName              = "Register"
	LoginRouteName                 = "Login"
	DashboardHomeRouteName         = "DashboardHome"
	DashboardProcessFlowsRouteName = "ProcessFlows"
	DashboardUserOrgsRouteName     = "DashboardUserOrgs"
	DashboardOrgHomeRouteName      = "DashboardOrgHome"
	DashboardUserHomeRouteName     = "DashboardUserHome"
	OrgSettingsUsersRouteName      = "OrgSettingsUsers"
	SingleDatabaseRouteName        = "SingleDatabase"
	SingleSystemRouteName          = "SingleSystem"
	SingleVendorRouteName          = "SingleVendor"
	SingleRiskRouteName            = "SingleRisk"
	SingleControlRouteName         = "SingleControl"
	SingleFlowRouteName            = "SingleFlow"
	SingleDocCatRouteName          = "SingleDocCat"
	SingleServerRouteName          = "SingleServer"
	SingleGLAccountRouteName       = "SingleGLAccount"
	FullGLAccountRouteName         = "FullGL"
	SingleDocRequestRouteName      = "SingleDocRequest"
	SingleSqlRequestRouteName      = "SingleSqlRequest"
	SingleScriptRequestRouteName   = "SingleScriptRequest"
	SingleDocumentationRouteName   = "SingleDocumentation"
	SingleClientDataRouteName      = "SingleClientData"
	SingleClientScriptRouteName    = "SingleClientScript"
	SingleGenericRequestRouteName  = "GenericRequest"
	SingleShellRouteName           = "SingleShell"
	// Api
	ApiRunCodeRouteName = "ApiRunCode"
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
	return core.EnvConfig.SelfUri + MustGetRouteUrl(r, params...)
}

func CreateOktaLoginUrl(idp string, state string, nonce string) string {
	envConfig := core.EnvConfig

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

func CreateOktaTokenUrl() string {
	return fmt.Sprintf("%s%s%s",
		core.EnvConfig.Okta.BaseUrl,
		core.EnvConfig.Login.AuthServerEndpoint,
		core.EnvConfig.Login.TokenEndpoint)
}

func CreateOktaKeyUrl() string {
	return fmt.Sprintf("%s%s%s?client_id=%s",
		core.EnvConfig.Okta.BaseUrl,
		core.EnvConfig.Login.AuthServerEndpoint,
		core.EnvConfig.Login.KeyEndpoint,
		core.EnvConfig.Login.ClientId)
}

func CreateOktaLogoutUrl(idToken string) string {
	envConfig := core.EnvConfig

	return fmt.Sprintf("%s%s%s?id_token_hint=%s&post_logout_redirect_uri=%s",
		envConfig.Okta.BaseUrl,
		envConfig.Login.AuthServerEndpoint,
		envConfig.Login.LogoutEndpoint,
		idToken,
		url.QueryEscape(MustGetRouteUrlAbsolute(LandingPageRouteName)))
}

func CreateOktaUserUpdateUrl(userId string) string {
	envConfig := core.EnvConfig

	return fmt.Sprintf("%s%s%s/%s",
		envConfig.Okta.BaseUrl,
		envConfig.Okta.ApiEndpoint,
		envConfig.Okta.UsersEndpoint,
		userId)
}
