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
	s.HandleFunc(core.DashboardHomeUrl, render.RenderDashboardHomePage).Methods("GET").Name(webcore.DashboardHomeRouteName)
	createOrganizationSubrouter(s)
	createUserSubrouter(s)
}

func createOrganizationSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardOrgUrl).Subrouter()
	s.Use(webcore.ObtainOrganizationInfoFromRequestInContextMiddleware)
	s.Use(webcore.CreateVerifyUserHasAccessToOrganizationMiddleware(func(w http.ResponseWriter, r *http.Request) {
		render.Render403(w, r)
	}))
	s.HandleFunc(core.DashboardOrgHomeUrl, render.RenderDashboardOrgHomePage).Methods("GET").Name(webcore.DashboardOrgHomeRouteName)
	s.PathPrefix(core.DashboardOrgFlowUrl).Handler(http.HandlerFunc(render.RenderDashboardProcessFlowsPage)).Methods("GET").Name(webcore.DashboardProcessFlowsRouteName)
	s.HandleFunc(core.DashboardOrgRiskEndpoint, render.RenderDashboardSingleRiskPage).Methods("GET").Name(webcore.DashboardSingleRiskRouteName)
	s.HandleFunc(core.DashboardOrgAllRiskEndpoint, render.RenderDashboardRisksPage).Methods("GET").Name(webcore.DashboardRisksRouteName)
}

func createUserSubrouter(r *mux.Router) {
	s := r.PathPrefix(core.DashboardUserUrl).Subrouter()
	s.Use(webcore.CreateVerifyUserHasAccessToUserMiddleware(func(w http.ResponseWriter, r *http.Request) {
		render.Render403(w, r)
	}))
	s.HandleFunc(core.DashboardUserHomeUrl, render.RenderDashboardUserHomePage).Methods("GET").Name(webcore.DashboardUserHomeRouteName)
}
