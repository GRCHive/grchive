package main

import (
	"github.com/gorilla/mux"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/render"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

func createDashboardSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardUrl).Subrouter()
	s.Use(webcore.CreateAuthenticatedRequestMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, webcore.MustGetRouteUrl(webcore.LoginRouteName), http.StatusTemporaryRedirect)
	}))
	s.Use(webcore.GrantAPIKeyMiddleware)

	s.HandleFunc(core.DashboardHomeUrl, render.RenderDashboardHomePage).Name(webcore.DashboardHomeRouteName)
	createOrganizationSubrouter(s)
	createUserSubrouter(s)
}

func createOrganizationSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardOrgUrl).Subrouter()
	s.Use(webcore.ObtainOrganizationInfoFromRequestInContextMiddleware)
	s.Use(webcore.CreateVerifyUserHasAccessToOrganizationMiddleware(func(w http.ResponseWriter, r *http.Request) {
		render.Render403(w, r)
	}))
	s.HandleFunc(core.DashboardOrgHomeUrl, render.RenderDashboardOrgHomePage).Name(webcore.DashboardOrgHomeRouteName)
	s.HandleFunc(core.DashboardOrgAllFlowsEndpoint, render.RenderDashboardProcessFlowsPage).Name(webcore.DashboardProcessFlowsRouteName)
	s.HandleFunc(core.DashboardOrgFlowEndpoint, render.RenderDashboardSingleFlowPage).Name(webcore.SingleFlowRouteName)
	s.HandleFunc(core.DashboardOrgRiskEndpoint, render.RenderDashboardSingleRiskPage).Name(webcore.SingleRiskRouteName)
	s.HandleFunc(core.DashboardOrgAllRiskEndpoint, render.RenderDashboardRisksPage)
	s.HandleFunc(core.DashboardOrgAllControlsEndpoint, render.RenderDashboardControlsPage)
	s.HandleFunc(core.DashboardOrgControlEndpoint, render.RenderDashboardSingleControlPage).Name(webcore.SingleControlRouteName)
	s.HandleFunc(core.DashboardOrgAllDocumentationEndpoint, render.RenderDocumentation)
	s.HandleFunc(core.DashboardOrgSingleDocCatEndpoint, render.RenderSingleDocCat).Name(webcore.SingleDocCatRouteName)
	s.HandleFunc(core.DashboardOrgAllDocRequestsEndpoint, render.RenderDocRequest)
	s.HandleFunc(core.DashboardOrgSingleDocRequestEndpoint, render.RenderSingleDocRequest).Name(webcore.SingleDocRequestRouteName)
	s.HandleFunc(core.DashboardOrgSingleSqlRequestEndpoint, render.RenderSingleSqlRequest).Name(webcore.SingleSqlRequestRouteName)
	s.HandleFunc(core.DashboardOrgSingleScriptRequestEndpoint, render.RenderSingleScriptRequest).Name(webcore.SingleScriptRequestRouteName)

	s.HandleFunc(core.DashboardOrgAllVendorsEndpoint, render.RenderVendors)
	s.HandleFunc(core.DashboardOrgSingleVendorEndpoint, render.RenderSingleVendor).Name(webcore.SingleVendorRouteName)
	s.HandleFunc(core.DashboardOrgSingleDocFileEndpoint, render.RenderSingleDocFile).Name(webcore.SingleDocumentationRouteName)
	s.HandleFunc(core.ApiAuditTrailPrefix, render.RenderAuditTrail)

	createOrganizationSettingsSubrouter(s)
	createOrganizationGLSubrouter(s)
	createOrganizationSystemSubrouter(s)
	createOrganizationAutomationSubrouter(s)
}

func createOrganizationGLSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardGeneralLedgerPrefix).Subrouter()
	s.HandleFunc(core.DashboardGeneralLedgerViewEndpoint, render.RenderDashboardGeneralLedger).Name(webcore.FullGLAccountRouteName)
	s.HandleFunc(core.DashboardOrgGLAccountEndpoint, render.RenderDashboardGLAccount).Name(webcore.SingleGLAccountRouteName)
}

func createOrganizationSystemSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardSystemsPrefix).Subrouter()
	s.HandleFunc(core.DashboardSystemHomeEndpoint, render.RenderSystemHome)
	s.HandleFunc(core.DashboardSingleSystemEndpoint, render.RenderSingleSystem).Name(webcore.SingleSystemRouteName)

	s.HandleFunc(core.DashboardDbSystemsEndpoint, render.RenderDbSystems)
	s.HandleFunc(core.DashboardSingleDbEndpoint, render.RenderSingleDb).Name(webcore.SingleDatabaseRouteName)

	s.HandleFunc(core.DashboardServersEndpoint, render.RenderServers)
	s.HandleFunc(core.DashboardSingleServerEndpoint, render.RenderSingleServer).Name(webcore.SingleServerRouteName)

	s.HandleFunc(core.DashboardShellEndpoint, render.RenderShells)
	s.HandleFunc(core.DashboardSingleShellEndpoint, render.RenderSingleShell).Name(webcore.SingleShellRouteName)
}

func createOrganizationSettingsSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardOrgSettingsPrefix).Subrouter()
	s.HandleFunc(core.DashboardOrgSettingsHomeEndpoint, render.RenderDashboardOrgSettingsHome)
	s.HandleFunc(core.DashboardOrgSettingsUsersEndpoint, render.RenderDashboardOrgSettingsUsers).Name(webcore.OrgSettingsUsersRouteName)
	s.HandleFunc(core.DashboardOrgSettingsRolesEndpoint, render.RenderDashboardOrgSettingsRoles)
	s.HandleFunc(core.DashboardOrgSettingsSingleRoleEndpoint, render.RenderDashboardOrgSettingsSingleRole)
}

func createOrganizationAutomationSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardAutomationPrefix).Subrouter()
	s.Use(webcore.CreateFeatureCheck(
		// Failure
		func(w http.ResponseWriter, r *http.Request) {
			render.Render404(w, r)
		},
		// Need Enable
		func(w http.ResponseWriter, r *http.Request) {
			render.RenderFeatureRequestPage(w, r, core.AutomationFeature, false)
		},
		// Pending
		func(w http.ResponseWriter, r *http.Request) {
			render.RenderFeatureRequestPage(w, r, core.AutomationFeature, true)
		},
		core.AutomationFeature,
	))

	s.HandleFunc("/schedule", render.RenderScriptSchedule)
	createOrganizationDataSubrouter(s)
	createOrganizationScriptSubrouter(s)
	createOrganizationLogsSubrouter(s)
}

func createOrganizationDataSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardDataPrefix).Subrouter()
	s.HandleFunc("/", render.RenderClientData)
	s.HandleFunc(core.DashboardSingleDataEndpoint, render.RenderSingleClientData).Name(webcore.SingleClientDataRouteName)
}

func createOrganizationScriptSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardScriptPrefix).Subrouter()
	s.HandleFunc("/", render.RenderClientScripts)
	s.HandleFunc(core.DashboardSingleScriptEndpoint, render.RenderSingleClientScript).Name(webcore.SingleClientScriptRouteName)
}

func createOrganizationLogsSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardLogsPrefix).Subrouter()
	s.HandleFunc("/", render.RenderLogs)
	s.HandleFunc(core.DashboardSingleBuildLogEndpoint, render.RenderSingleBuildLog)
	s.HandleFunc(core.DashboardSingleScriptRunLogEndpoint, render.RenderSingleRunLog)
}

func createUserSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardUserPrefix).Subrouter()
	s.Use(webcore.ObtainUserInfoFromRequestInContextMiddleware)
	s.Use(webcore.CreateVerifyUserHasAccessToUserMiddleware(func(w http.ResponseWriter, r *http.Request) {
		render.Render403(w, r)
	}))
	s.HandleFunc(core.DashboardUserHomeUrl, render.RenderDashboardUserHomePage).Name(webcore.DashboardUserHomeRouteName)
	s.HandleFunc(core.DashboardUserOrgUrl, render.RenderDashboardUserOrgsPage).Name(webcore.DashboardUserOrgsRouteName)
	s.HandleFunc(core.DashboardUserProfileUrl, render.RenderDashboardUserProfilePage)
	s.HandleFunc(core.DashboardUserNotificationsUrl, render.RenderDashboardUserNotificationsPage)
}
