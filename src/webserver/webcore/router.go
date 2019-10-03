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
	GettingStartedRouteName     RouteName = "GettingStarted"
	ContactUsRouteName                    = "ContactUs"
	LandingPageRouteName                  = "LandingPage"
	LoginRouteName                        = "Login"
	LearnMoreRouteName                    = "LearnMore"
	GettingStartedPostRouteName           = "GettingStartedPost"
	LoginPostRouteName                    = "LoginPost"
	SamlCallbackRouteName                 = "SamlCallback"
	DashboardHomeRouteName                = "DashboardHome"
	DashboardOrgHomeRouteName             = "DashboardOrgHome"
)

var globalRouter *mux.Router

func RegisterRouter(r *mux.Router) {
	globalRouter = r
}

func MustGetRouteUrl(r RouteName) string {
	route := globalRouter.Get(string(r))
	if route == nil {
		core.Warning("No Route: " + string(r))
		return "/404"
	}

	url, err := route.URL()
	if err != nil {
		core.Warning("Bad Route: " + string(r))
		return "/404"
	}
	return url.String()
}

func MustGetRouteUrlAbsolute(r RouteName) string {
	return core.LoadEnvConfig().SelfUri + MustGetRouteUrl(r)
}

func CreateOktaLoginUrl(idp string, state string, nonce string) string {
	envConfig := core.LoadEnvConfig()

	return fmt.Sprintf("%s%s?idp=%s&client_id=%s&response_type=%s&response_mode=%s&scope=%s&redirect_uri=%s&state=%s&=nonce=%s",
		envConfig.Login.BaseUrl,
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

var OktaTokenUrl string = fmt.Sprintf("%s%s",
	core.LoadEnvConfig().Login.BaseUrl,
	core.LoadEnvConfig().Login.TokenEndpoint)

var OktaKeyUrl string = fmt.Sprintf("%s%s?client_id=%s",
	core.LoadEnvConfig().Login.BaseUrl,
	core.LoadEnvConfig().Login.KeyEndpoint,
	core.LoadEnvConfig().Login.ClientId)
