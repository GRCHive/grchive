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
	registerAPIv2ServerPaths(s)
}

func registerAPIv2ServerPaths(r *mux.Router) {
	s := r.PathPrefix("/server").Subrouter()

	ss := s.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgServerQueryId)).Subrouter()
	ss.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgServerQueryId))

	ssc := ss.PathPrefix("/connection").Subrouter()
	ssc.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			allServerConnections,
			core.ResourceAccessBundle{core.ResourceServers, core.AccessView},
		),
	).Methods("GET")

	sscsp := ssc.PathPrefix("/ssh/password").Subrouter()
	sscsp.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			newServerConnectionSSHPassword,
			core.ResourceAccessBundle{core.ResourceServers, core.AccessEdit},
		),
	).Methods("POST")

	sscspc := sscsp.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgServerSshConnectionQueryId)).Subrouter()
	sscspc.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgServerSshConnectionQueryId))
	sscspc.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			deleteServerConnectionSSHPassword,
			core.ResourceAccessBundle{core.ResourceServers, core.AccessEdit},
		),
	).Methods("DELETE")

	sscspc.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			getServerConnectionSSHPassword,
			core.ResourceAccessBundle{core.ResourceServers, core.AccessView},
		),
	).Methods("GET")

	sscspc.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			editServerConnectionSSHPassword,
			core.ResourceAccessBundle{core.ResourceServers, core.AccessView},
		),
	).Methods("PUT")

	sscsk := ssc.PathPrefix("/ssh/key").Subrouter()
	sscsk.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			newServerConnectionSSHKey,
			core.ResourceAccessBundle{core.ResourceServers, core.AccessEdit},
		),
	).Methods("POST")

	sscskc := sscsk.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgServerSshKeyQueryId)).Subrouter()
	sscskc.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgServerSshKeyQueryId))
	sscskc.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			deleteServerConnectionSSHKey,
			core.ResourceAccessBundle{core.ResourceServers, core.AccessEdit},
		),
	).Methods("DELETE")

	sscskc.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			getServerConnectionSSHKey,
			core.ResourceAccessBundle{core.ResourceServers, core.AccessView},
		),
	).Methods("GET")

	sscskc.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			editServerConnectionSSHKey,
			core.ResourceAccessBundle{core.ResourceServers, core.AccessView},
		),
	).Methods("PUT")

}

func registerAPIv2ShellPaths(r *mux.Router) {
	// All organization shell scripts
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

	sr := s.PathPrefix("/run").Subrouter()
	sr.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			allShellRuns,
			core.ResourceAccessBundle{core.ResourceShellRun, core.AccessView},
		),
	).Methods("GET")

	srr := sr.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgShellRunQueryId)).Subrouter()
	srr.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgShellRunQueryId))

	srr.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			getShellRun,
			core.ResourceAccessBundle{core.ResourceShellRun, core.AccessView},
			core.ResourceAccessBundle{core.ResourceShell, core.AccessView},
		),
	).Methods("GET")

	// Individual shell scripts
	ss := s.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgShellScriptQueryId)).Subrouter()
	ss.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgShellScriptQueryId))
	ss.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			deleteShellScript,
			core.ResourceAccessBundle{core.ResourceShell, core.AccessManage},
		),
	).Methods("DELETE")

	ss.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			getShellScript,
			core.ResourceAccessBundle{core.ResourceShell, core.AccessView},
		),
	).Methods("GET")

	ss.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			editShellScript,
			core.ResourceAccessBundle{core.ResourceShell, core.AccessEdit},
		),
	).Methods("PUT")

	// Shell script versions
	ssv := ss.PathPrefix("/version").Subrouter()
	ssv.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			allShellVersions,
			core.ResourceAccessBundle{core.ResourceShell, core.AccessView},
		),
	).Methods("GET")

	ssv.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			newShellVersion,
			core.ResourceAccessBundle{core.ResourceShell, core.AccessManage},
		),
	).Methods("POST")

	ssvv := ssv.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgShellScriptVersionQueryId)).Subrouter()
	ssvv.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgShellScriptVersionQueryId))
	ssvv.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			getShellVersion,
			core.ResourceAccessBundle{core.ResourceShell, core.AccessView},
		),
	).Methods("GET")

	ssvvr := ssvv.PathPrefix("/run").Subrouter()
	ssvvr.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			runShellVersion,
			core.ResourceAccessBundle{core.ResourceShellRun, core.AccessManage},
		),
	).Methods("POST")
}

func registerAPIv2RequestsPaths(r *mux.Router) {
	s := r.PathPrefix("/requests").Subrouter()

	scriptRouter := s.PathPrefix("/scripts").Subrouter()
	scriptRouter.HandleFunc("/", allGenericRequestsScripts).Methods("GET")

	singleScriptRouter := scriptRouter.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgScriptRequestQueryId)).Subrouter()
	singleScriptRouter.Use(webcore.CreateObtainGenericRequestInContext(core.DashboardOrgScriptRequestQueryId))
	singleScriptRouter.HandleFunc("/", getGenericRequestScript).Methods("GET")
	singleScriptRouter.HandleFunc("/approval", approveDenyScriptRunRequest).Methods("POST")

	shellRouter := s.PathPrefix("/shell").Subrouter()
	shellRouter.HandleFunc("/", allGenericRequestsShellScripts).Methods("GET")

	singleShellRouter := shellRouter.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgShellRequestQueryId)).Subrouter()
	singleShellRouter.Use(webcore.CreateObtainGenericRequestInContext(core.DashboardOrgShellRequestQueryId))
	singleShellRouter.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			getGenericRequestShell,
			core.ResourceAccessBundle{core.ResourceShell, core.AccessView},
		),
	).Methods("GET")
	singleShellRouter.HandleFunc(
		"/approval",
		webcore.CreateACLCheckPermissionHandler(
			approveDenyShellRunRequest,
			core.ResourceAccessBundle{core.ResourceShellRun, core.AccessEdit},
		),
	).Methods("POST")

	singleReqRouter := s.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgScriptRequestQueryId)).Subrouter()
	singleReqRouter.Use(webcore.CreateObtainGenericRequestInContext(core.DashboardOrgScriptRequestQueryId))
	singleReqRouter.HandleFunc("/", getGenericRequest).Methods("GET")
	singleReqRouter.HandleFunc("/", editGenericRequest).Methods("PUT")
	singleReqRouter.HandleFunc("/", deleteGenericRequest).Methods("DELETE")
	singleReqRouter.HandleFunc("/approval", getGenericApproval).Methods("GET")
}
