package main

import (
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/render"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
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
	s.HandleFunc(core.DashboardOrgFlowEndpoint, render.RenderDashboardSingleFlowPage).Name(webcore.DashboardSingleFlowRouteName)
	s.HandleFunc(core.DashboardOrgRiskEndpoint, render.RenderDashboardSingleRiskPage).Name(webcore.DashboardSingleRiskRouteName)
	s.HandleFunc(core.DashboardOrgAllRiskEndpoint, render.RenderDashboardRisksPage).Name(webcore.DashboardRisksRouteName)
	s.HandleFunc(core.DashboardOrgAllControlsEndpoint, render.RenderDashboardControlsPage).Name(webcore.DashboardControlsRouteName)
	s.HandleFunc(core.DashboardOrgControlEndpoint, render.RenderDashboardSingleControlPage).Name(webcore.DashboardSingleControlRouteName)
	s.HandleFunc(core.DashboardOrgAllDocumentationEndpoint, render.RenderDocumentation).Name(webcore.DashboardDocumentationsRouteName)
	s.HandleFunc(core.DashboardOrgSingleDocCatEndpoint, render.RenderSingleDocCat).Name(webcore.DashboardSingleDocumentationRouteName)

	createOrganizationSettingsSubrouter(s)
	createOrganizationGLSubrouter(s)
	createOrganizationSystemSubrouter(s)
}

func createOrganizationGLSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardGeneralLedgerPrefix).Subrouter()
	s.HandleFunc(core.DashboardGeneralLedgerViewEndpoint, render.RenderDashboardGeneralLedger).Name(webcore.GeneralLedgerRouteName)
	s.HandleFunc(core.DashboardOrgGLAccountEndpoint, render.RenderDashboardGLAccount).Name(webcore.GLAccountRouteName)
}

func createOrganizationSystemSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardSystemsPrefix).Subrouter()
	s.HandleFunc(core.DashboardSystemHomeEndpoint, render.RenderSystemHome).Name(webcore.SystemHomeRouteName)
	s.HandleFunc(core.DashboardSingleSystemEndpoint, render.RenderSingleSystem).Name(webcore.SingleSystemRouteName)

	s.HandleFunc(core.DashboardDbSystemsEndpoint, render.RenderDbSystems).Name(webcore.DbSystemsRouteName)
	s.HandleFunc(core.DashboardSingleDbEndpoint, render.RenderSingleDb).Name(webcore.SingleDbRouteName)

	s.HandleFunc(core.DashboardInfraSystemsEndpoint, render.RenderInfraSystems).Name(webcore.InfraSystemsRouteName)
}

func createOrganizationSettingsSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardOrgSettingsPrefix).Subrouter()
	s.HandleFunc(core.DashboardOrgSettingsHomeEndpoint, render.RenderDashboardOrgSettingsHome).Name(webcore.OrgSettingsHomeRouteName)
	s.HandleFunc(core.DashboardOrgSettingsUsersEndpoint, render.RenderDashboardOrgSettingsUsers).Name(webcore.OrgSettingsUsersRouteName)
	s.HandleFunc(core.DashboardOrgSettingsRolesEndpoint, render.RenderDashboardOrgSettingsRoles).Name(webcore.OrgSettingsRolesRouteName)
	s.HandleFunc(core.DashboardOrgSettingsSingleRoleEndpoint, render.RenderDashboardOrgSettingsSingleRole).Name(webcore.OrgSettingsSingleRoleRouteName)
}

func createUserSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardUserPrefix).Subrouter()
	s.Use(webcore.ObtainUserInfoFromRequestInContextMiddleware)
	s.Use(webcore.CreateVerifyUserHasAccessToUserMiddleware(func(w http.ResponseWriter, r *http.Request) {
		render.Render403(w, r)
	}))
	s.HandleFunc(core.DashboardUserHomeUrl, render.RenderDashboardUserHomePage).Name(webcore.DashboardUserHomeRouteName)
	s.HandleFunc(core.DashboardUserOrgUrl, render.RenderDashboardUserOrgsPage).Name(webcore.DashboardUserOrgsRouteName)
	s.HandleFunc(core.DashboardUserProfileUrl, render.RenderDashboardUserProfilePage).Name(webcore.DashboardUserProfileRouteName)
}
