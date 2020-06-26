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
	registerAPIv2GenericRequestsPaths(s)
	registerAPIv2ServerPaths(s)
	registerAPIv2DatabasePaths(s)
	registerAPIv2SystemPaths(s)
	registerAPIv2IntegrationPaths(s)
	registerAPIv2PBCRequestsPaths(s)
	registerAPIv2AnalyticsPaths(s)
}

func registerAPIv2DatabasePaths(r *mux.Router) {
	s := r.PathPrefix("/database").Subrouter()

	sd := s.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgDbQueryId)).Subrouter()
	sd.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgDbQueryId))

	sds := sd.PathPrefix("/settings").Subrouter()
	sds.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			getDatabaseSettings,
			core.ResourceAccessBundle{core.ResourceDatabases, core.AccessView},
		),
	).Methods("GET")

	sds.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			editDatabaseSettings,
			core.ResourceAccessBundle{core.ResourceDatabases, core.AccessEdit},
		),
	).Methods("PUT")
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

func registerAPIv2PBCRequestsPaths(r *mux.Router) {
	s := r.PathPrefix("/pbc").Subrouter()

	singleReq := s.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgDocRequestQueryId)).Subrouter()
	singleReq.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgDocRequestQueryId))

	singleReqFiles := singleReq.PathPrefix("/files").Subrouter()
	singleReqFiles.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			newDocRequestFileLinks,
			core.ResourceAccessBundle{core.ResourceDocRequests, core.AccessEdit},
			core.ResourceAccessBundle{core.ResourceControlDocumentationMetadata, core.AccessEdit},
		),
	).Methods("POST")

	singleReqControl := singleReq.PathPrefix("/control").Subrouter()
	singleReqSingleControl := singleReqControl.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgControlQueryId)).Subrouter()
	singleReqSingleControl.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgControlQueryId))

	singleReqSingleControl.HandleFunc("/folders",
		webcore.CreateACLCheckPermissionHandler(
			allDocRequestControlFolderLinks,
			core.ResourceAccessBundle{core.ResourceControls, core.AccessView},
			core.ResourceAccessBundle{core.ResourceDocRequests, core.AccessView},
		),
	).Methods("GET")
}

func registerAPIv2GenericRequestsPaths(r *mux.Router) {
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

	singleReqRouter := s.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgGenericRequestQueryId)).Subrouter()
	singleReqRouter.Use(webcore.CreateObtainGenericRequestInContext(core.DashboardOrgGenericRequestQueryId))
	singleReqRouter.HandleFunc("/", getGenericRequest).Methods("GET")
	singleReqRouter.HandleFunc("/", editGenericRequest).Methods("PUT")
	singleReqRouter.HandleFunc("/", deleteGenericRequest).Methods("DELETE")
	singleReqRouter.HandleFunc("/approval", getGenericApproval).Methods("GET")
}

func registerAPIv2SystemPaths(r *mux.Router) {
	s := r.PathPrefix("/system").Subrouter()

	singleSystemRouter := s.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgSystemQueryId)).Subrouter()
	singleSystemRouter.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgSystemQueryId))
	registerAPIv2SystemIntegrationPaths(singleSystemRouter)
}

// This function currently assumes that the only time we deal with integrations is under the context of a system.
func registerAPIv2SystemIntegrationPaths(r *mux.Router) {
	s := r.PathPrefix("/integration").Subrouter()
	s.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			allIntegrations,
			core.ResourceAccessBundle{core.ResourceSystems, core.AccessView},
			core.ResourceAccessBundle{core.ResourceIntegrations, core.AccessView},
		),
	)

	sapErpRouter := s.PathPrefix("/sap/erp").Subrouter()
	sapErpRouter.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			newSapErpIntegration,
			core.ResourceAccessBundle{core.ResourceSystems, core.AccessEdit},
			core.ResourceAccessBundle{core.ResourceIntegrations, core.AccessManage},
			core.ResourceAccessBundle{core.ResourceSapErp, core.AccessManage},
		),
	).Methods("POST")
}

func registerAPIv2IntegrationPaths(r *mux.Router) {
	s := r.PathPrefix("/integration").Subrouter()

	ss := s.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgIntegrationQueryId)).Subrouter()
	ss.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgIntegrationQueryId))

	ss.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			editGenericIntegration,
			core.ResourceAccessBundle{core.ResourceIntegrations, core.AccessEdit},
		),
	).Methods("PUT")

	ss.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			deleteGenericIntegration,
			core.ResourceAccessBundle{core.ResourceIntegrations, core.AccessManage},
		),
	).Methods("DELETE")

	ssSapErp := ss.PathPrefix("/sap/erp").Subrouter()
	ssSapErp.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			getSapErpIntegration,
			core.ResourceAccessBundle{core.ResourceSapErp, core.AccessView},
		),
	).Methods("GET")

	ssSapErp.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			editSapErpIntegration,
			core.ResourceAccessBundle{core.ResourceSapErp, core.AccessEdit},
		),
	).Methods("PUT")

	ssSapErpRfc := ssSapErp.PathPrefix("/rfc").Subrouter()
	ssSapErpRfc.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			allSapErpRfc,
			core.ResourceAccessBundle{core.ResourceSapErp, core.AccessView},
		),
	).Methods("GET")

	ssSapErpRfc.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			newSapErpRfc,
			core.ResourceAccessBundle{core.ResourceSapErp, core.AccessEdit},
		),
	).Methods("POST")

	ssSapErpRfcSingle := ssSapErpRfc.PathPrefix(fmt.Sprintf("/{%s}", core.DashboardOrgSapErpRfcQueryId)).Subrouter()
	ssSapErpRfcSingle.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgSapErpRfcQueryId))

	ssSapErpRfcSingle.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			deleteSapErpRfc,
			core.ResourceAccessBundle{core.ResourceSapErp, core.AccessEdit},
		),
	).Methods("DELETE")

	ssSapErpRfcSingleVersions := ssSapErpRfcSingle.PathPrefix("/version").Subrouter()
	ssSapErpRfcSingleVersions.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			allSapErpRfcVersions,
			core.ResourceAccessBundle{core.ResourceSapErp, core.AccessView},
		),
	).Methods("GET")

	ssSapErpRfcSingleVersions.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			newSapErpRfcVersion,
			core.ResourceAccessBundle{core.ResourceSapErp, core.AccessEdit},
		),
	).Methods("POST")

	ssSapErpRfcSingleVersionSingle := ssSapErpRfcSingleVersions.PathPrefix(
		fmt.Sprintf("/{%s}", core.DashboardOrgSapErpRfcVersionQueryId),
	).Subrouter()
	ssSapErpRfcSingleVersionSingle.Use(webcore.CreateObtainResourceInContextMiddleware(core.DashboardOrgSapErpRfcVersionQueryId))
	ssSapErpRfcSingleVersionSingle.HandleFunc(
		"/",
		webcore.CreateACLCheckPermissionHandler(
			getSapErpRfcVersion,
			core.ResourceAccessBundle{core.ResourceSapErp, core.AccessView},
		),
	).Methods("GET")
}

func registerAPIv2AnalyticsPaths(r *mux.Router) {
	s := r.PathPrefix("/analytics").Subrouter()
	registerAPIv2PbcAnalytics(s)
}

func registerAPIv2PbcAnalytics(r *mux.Router) {
	s := r.PathPrefix("/pbc").Subrouter()
	s.HandleFunc(
		"/overall",
		webcore.CreateACLCheckPermissionHandler(
			getOverallPbcAnalytics,
			core.ResourceAccessBundle{core.ResourceDocRequests, core.AccessView},
		),
	)

	s.HandleFunc(
		"/category/requester",
		webcore.CreateACLCheckPermissionHandler(
			getCategoryRequesterPbcAnalytics,
			core.ResourceAccessBundle{core.ResourceDocRequests, core.AccessView},
		),
	)

	s.HandleFunc(
		"/category/assignee",
		webcore.CreateACLCheckPermissionHandler(
			getCategoryAssigneePbcAnalytics,
			core.ResourceAccessBundle{core.ResourceDocRequests, core.AccessView},
		),
	)

	s.HandleFunc(
		"/category/cat",
		webcore.CreateACLCheckPermissionHandler(
			getCategoryDocCatPbcAnalytics,
			core.ResourceAccessBundle{core.ResourceDocRequests, core.AccessView},
		),
	)

	s.HandleFunc(
		"/category/flow",
		webcore.CreateACLCheckPermissionHandler(
			getCategoryProcessFlowPbcAnalytics,
			core.ResourceAccessBundle{core.ResourceDocRequests, core.AccessView},
		),
	)

	s.HandleFunc(
		"/category/control",
		webcore.CreateACLCheckPermissionHandler(
			getCategoryControlPbcAnalytics,
			core.ResourceAccessBundle{core.ResourceDocRequests, core.AccessView},
		),
	)

	s.HandleFunc(
		"/category/risk",
		webcore.CreateACLCheckPermissionHandler(
			getCategoryRiskPbcAnalytics,
			core.ResourceAccessBundle{core.ResourceDocRequests, core.AccessView},
		),
	)

	s.HandleFunc(
		"/category/gl",
		webcore.CreateACLCheckPermissionHandler(
			getCategoryGLPbcAnalytics,
			core.ResourceAccessBundle{core.ResourceDocRequests, core.AccessView},
		),
	)

	s.HandleFunc(
		"/category/system",
		webcore.CreateACLCheckPermissionHandler(
			getCategorySystemPbcAnalytics,
			core.ResourceAccessBundle{core.ResourceDocRequests, core.AccessView},
		),
	)
}
