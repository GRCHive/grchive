package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/webcore"
)

func registerAPIv2Paths(r *mux.Router) {
	s := r.PathPrefix("/v2").Subrouter()
	registerAPIv2Org(s)
}

func registerAPIv2Org(r *mux.Router) {
	s := r.PathPrefix(fmt.Sprintf("/org/{%s}", core.DashboardOrgOrgQueryId)).Subrouter()
	s.Use(webcore.ObtainOrganizationFromIdInRequestInContextMiddleware)
	s.Use(webcore.ObtainRoleFromRequestInContextMiddleware)

	registerAPIv2ShellPaths(s)
	registerAPIv2RequestsPaths(s)
}

func registerAPIv2ShellPaths(r *mux.Router) {
	s := r.PathPrefix("/shell").Subrouter()
	s.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			allShells,
			core.ResourceAccessBundle{core.ResourceShell, core.AccessView},
		),
	).Methods("GET")

	s.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			newShell,
			core.ResourceAccessBundle{core.ResourceShell, core.AccessManage},
		),
	).Methods("POST")

	sv := s.PathPrefix("/version").Subrouter()
	sv.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			allShellVersions,
			core.ResourceAccessBundle{core.ResourceShell, core.AccessView},
		),
	).Methods("GET")

	ss := s.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgShellScriptQueryId)).Subrouter()
	ss.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgShellScriptQueryId))
	ss.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			deleteShellScript,
			core.ResourceAccessBundle{core.ResourceShell, core.AccessManage},
		),
	).Methods("DELETE")
}

func registerAPIv2RequestsPaths(r *mux.Router) {
	s := r.PathPrefix("/requests").Subrouter()

	scriptRouter := s.PathPrefix("/scripts").Subrouter()
	scriptRouter.HandleFunc("/", allGenericRequestsScripts).Methods("GET")

	singleScriptRouter := scriptRouter.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgScriptRequestQueryId)).Subrouter()
	singleScriptRouter.Use(webcore.CreateObtainGenericRequestInContext(core.DashboardOrgScriptRequestQueryId))
	singleScriptRouter.HandleFunc("/", getGenericRequestScript).Methods("GET")
	singleScriptRouter.HandleFunc("/approval", approveDenyScriptRunRequest).Methods("POST")

	singleReqRouter := s.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgScriptRequestQueryId)).Subrouter()
	singleReqRouter.Use(webcore.CreateObtainGenericRequestInContext(core.DashboardOrgScriptRequestQueryId))
	singleReqRouter.HandleFunc("/", getGenericRequest).Methods("GET")
	singleReqRouter.HandleFunc("/", editGenericRequest).Methods("PUT")
	singleReqRouter.HandleFunc("/", deleteGenericRequest).Methods("DELETE")
	singleReqRouter.HandleFunc("/approval", getGenericApproval).Methods("GET")
}
