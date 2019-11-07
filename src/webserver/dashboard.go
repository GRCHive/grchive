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
}

func createUserSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardUserUrl).Subrouter()
	s.Use(webcore.CreateVerifyUserHasAccessToUserMiddleware(func(w http.ResponseWriter, r *http.Request) {
		render.Render403(w, r)
	}))
	s.HandleFunc(core.DashboardUserHomeUrl, render.RenderDashboardUserHomePage).Name(webcore.DashboardUserHomeRouteName)
}
