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
	RegisterRouteName                              = "Register"
	LogoutRouteName                                = "Logout"
	LearnMoreRouteName                             = "LearnMore"
	GettingStartedPostRouteName                    = "GettingStartedPost"
	LoginPostRouteName                             = "LoginPost"
	RegisterPostRouteName                          = "RegisterPost"
	SamlCallbackRouteName                          = "SamlCallback"
	DashboardHomeRouteName                         = "DashboardHome"
	DashboardOrgHomeRouteName                      = "DashboardOrgHome"
	DashboardUserHomeRouteName                     = "DashboardUserHome"
	DashboardUserOrgsRouteName                     = "DashboardUserOrgs"
	DashboardUserProfileRouteName                  = "DashboardUserProfile"
	UserProfileEditRouteName                       = "UserProfilePost"
	UserGetOrgsRouteName                           = "UserOrgsGet"
	DashboardProcessFlowsRouteName                 = "ProcessFlows"
	DashboardRisksRouteName                        = "Risks"
	DashboardSingleRiskRouteName                   = "SingleRisk"
	DashboardSingleControlRouteName                = "SingleControl"
	DashboardSingleFlowRouteName                   = "SingleFlow"
	DashboardControlsRouteName                     = "Controls"
	NewProcessFlowRouteName                        = "NewProcessFlow"
	DeleteProcessFlowRouteName                     = "DeleteProcessFlow"
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
	EditRiskRouteName                              = "EditRisk"
	AddRiskToNodeRouteName                         = "AddRiskToNodeRisk"
	GetAllRiskRouteName                            = "GetAllRisk"
	GetSingleRiskRouteName                         = "GetSingleRisk"
	GetAllOrgUsersRouteName                        = "GetAllOrgUsers"
	NewControlRouteName                            = "NewControl"
	ControlTypesRouteName                          = "ControlTypes"
	DeleteControlRouteName                         = "DeleteControl"
	AddControlRouteName                            = "AddControl"
	EditControlRouteName                           = "EditControl"
	GetAllControlRouteName                         = "GetAllControl"
	GetSingleControlRouteName                      = "GetSingleControl"
	NewControlDocCatRouteName                      = "NewControlDocCat"
	EditControlDocCatRouteName                     = "EditControlDocCat"
	DeleteControlDocCatRouteName                   = "DeleteControlDocCat"
	UploadControlDocRouteName                      = "UploadControlDoc"
	GetControlDocRouteName                         = "GetControlDoc"
	DeleteControlDocRouteName                      = "DeleteControlDoc"
	DownloadControlDocRouteName                    = "DownloadControlDoc"
	EmailVerifyRouteName                           = "VerifyEmail"
	ResendVerificationRouteName                    = "ResendVerification"
	OrgSettingsHomeRouteName                       = "OrgSettingsHome"
	OrgSettingsUsersRouteName                      = "OrgSettingsUsers"
	OrgSettingsRolesRouteName                      = "OrgSettingsRoles"
	SendInviteRouteName                            = "SendInvite"
	AcceptInviteRouteName                          = "AcceptInvite"
	GetOrgRolesRouteName                           = "GetOrgRoles"
	GetSingleRoleRouteName                         = "GetSingleRole"
	NewRoleRouteName                               = "NewRole"
	EditRoleRouteName                              = "EditRole"
	DeleteRoleRouteName                            = "DeleteRole"
	AddUsersToRoleRouteName                        = "AddUsersToRole"
	OrgSettingsSingleRoleRouteName                 = "OrgSettingsSingleRole"
	GeneralLedgerRouteName                         = "GeneralLedger"
	GLAccountRouteName                             = "GLAccount"
	ApiGetGLRouteName                              = "ApiGetGL"
	ApiNewGLCatRouteName                           = "ApiNewGLCat"
	ApiEditGLCatRouteName                          = "ApiEditGLCat"
	ApiDeleteGLCatRouteName                        = "ApiDeleteGLCat"
	ApiNewGLAccRouteName                           = "ApiNewGLAcc"
	ApiGetGLAccRouteName                           = "ApiGetGLAcc"
	ApiEditGLAccRouteName                          = "ApiEditGLAcc"
	ApiDeleteGLAccRouteName                        = "ApiDeleteGLAcc"
	SystemHomeRouteName                            = "SystemHome"
	DbSystemsRouteName                             = "DbSystesm"
	InfraSystemsRouteName                          = "InfraSystems"
	ApiNewSystemRouteName                          = "ApiNewSystem"
	ApiSystemAllRouteName                          = "ApiSystemAll"
	ApiGetSystemRouteName                          = "ApiGetSystem"
	ApiEditSystemRouteName                         = "ApiEditSystem"
	ApiDeleteSystemRouteName                       = "ApiDeleteSystem"
	SingleSystemRouteName                          = "SingleSystem"
	SingleDbRouteName                              = "SingleDb"
	ApiNewDbRouteName                              = "ApiNewDb"
	ApiAllDbRouteName                              = "ApiAllDb"
	ApiTypesDbRouteName                            = "ApiTypesDb"
	ApiGetDbRouteName                              = "ApiGetDb"
	ApiEditDbRouteName                             = "ApiEditDb"
	ApiDeleteDbRouteName                           = "ApiDeleteDb"
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
