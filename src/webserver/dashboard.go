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
	s.HandleFunc(core.DashboardOrgFlowEndpoint, render.RenderDashboardSingleFlowPage)
	s.HandleFunc(core.DashboardOrgRiskEndpoint, render.RenderDashboardSingleRiskPage)
	s.HandleFunc(core.DashboardOrgAllRiskEndpoint, render.RenderDashboardRisksPage)
	s.HandleFunc(core.DashboardOrgAllControlsEndpoint, render.RenderDashboardControlsPage)
	s.HandleFunc(core.DashboardOrgControlEndpoint, render.RenderDashboardSingleControlPage)
	s.HandleFunc(core.DashboardOrgAllDocumentationEndpoint, render.RenderDocumentation)
	s.HandleFunc(core.DashboardOrgSingleDocCatEndpoint, render.RenderSingleDocCat)
	s.HandleFunc(core.DashboardOrgAllDocRequestsEndpoint, render.RenderDocRequest)
	s.HandleFunc(core.DashboardOrgSingleDocRequestEndpoint, render.RenderSingleDocRequest)
	s.HandleFunc(core.DashboardOrgSingleSqlRequestEndpoint, render.RenderSingleSqlRequest)
	s.HandleFunc(core.DashboardOrgAllVendorsEndpoint, render.RenderVendors)
	s.HandleFunc(core.DashboardOrgSingleVendorEndpoint, render.RenderSingleVendor)
	s.HandleFunc(core.DashboardOrgSingleDocFileEndpoint, render.RenderSingleDocFile)

	createOrganizationSettingsSubrouter(s)
	createOrganizationGLSubrouter(s)
	createOrganizationSystemSubrouter(s)
}

func createOrganizationGLSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardGeneralLedgerPrefix).Subrouter()
	s.HandleFunc(core.DashboardGeneralLedgerViewEndpoint, render.RenderDashboardGeneralLedger)
	s.HandleFunc(core.DashboardOrgGLAccountEndpoint, render.RenderDashboardGLAccount)
}

func createOrganizationSystemSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardSystemsPrefix).Subrouter()
	s.HandleFunc(core.DashboardSystemHomeEndpoint, render.RenderSystemHome)
	s.HandleFunc(core.DashboardSingleSystemEndpoint, render.RenderSingleSystem)

	s.HandleFunc(core.DashboardDbSystemsEndpoint, render.RenderDbSystems)
	s.HandleFunc(core.DashboardSingleDbEndpoint, render.RenderSingleDb)

	s.HandleFunc(core.DashboardServersEndpoint, render.RenderServers)
	s.HandleFunc(core.DashboardSingleServerEndpoint, render.RenderSingleServer)
}

func createOrganizationSettingsSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardOrgSettingsPrefix).Subrouter()
	s.HandleFunc(core.DashboardOrgSettingsHomeEndpoint, render.RenderDashboardOrgSettingsHome)
	s.HandleFunc(core.DashboardOrgSettingsUsersEndpoint, render.RenderDashboardOrgSettingsUsers).Name(webcore.OrgSettingsUsersRouteName)
	s.HandleFunc(core.DashboardOrgSettingsRolesEndpoint, render.RenderDashboardOrgSettingsRoles)
	s.HandleFunc(core.DashboardOrgSettingsSingleRoleEndpoint, render.RenderDashboardOrgSettingsSingleRole)
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
}
