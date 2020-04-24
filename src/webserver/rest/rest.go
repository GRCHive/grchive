package rest

import (
	"github.com/gorilla/mux"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/webcore"
)

func RegisterPaths(r *mux.Router) {
	altApiRouter := r.NewRoute().Subrouter()
	altApiRouter.Use(webcore.ObtainUserSessionInContextMiddleware)

	altApiRouter.HandleFunc(core.GetStartedUrl, postGettingStartedInterest).Methods("POST")

	// core.LoginURL is the POST request the user will send with just their email.
	// core.CreateSamlCallbackUrl() is the GET request the user's SAML IdP will redirect to upon
	// successful login.
	altApiRouter.HandleFunc(core.LoginUrl, postLogin).Methods("POST")
	altApiRouter.HandleFunc(core.RegisterUrl, postRegister).Methods("POST")
	altApiRouter.HandleFunc(core.CreateSamlCallbackUrl(), getSamlLoginCallback).Methods("GET").Name(webcore.SamlCallbackRouteName)
	altApiRouter.HandleFunc(core.LogoutUrl, getLogout).Methods("GET")
	altApiRouter.HandleFunc(core.VerifyEmailUrl, verifyUserEmail).Methods("GET").Name(webcore.EmailVerifyRouteName)
	altApiRouter.HandleFunc(core.AcceptInviteUrl, acceptInviteToOrganization).Methods("GET").Name(webcore.AcceptInviteRouteName)

	// REST API
	registerAPIPaths(r)
}

func registerAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiUrl).Subrouter()
	s.Use(webcore.ObtainAPIKeyRoleInContextMiddleware)

	registerAuditTrailAPIPaths(s)
	registerInviteAPIPaths(s)
	registerVerificationAPIPaths(s)
	registerUserAPIPaths(s)
	registerOrgAPIPaths(s)
	registerProcessFlowAPIPaths(s)
	registerProcessFlowNodesAPIPaths(s)
	registerProcessFlowIOAPIPaths(s)
	registerProcessFlowEdgesAPIPaths(s)
	registerRiskAPIPaths(s)
	registerControlAPIPaths(s)
	registerControlDocumentationAPIPaths(s)
	registerRoleAPIPaths(s)
	registerGeneralLedgerAPIPaths(s)
	registerITAPIPaths(s)
	registerDocRequestsAPIPaths(s)
	registerCommentsAPIPaths(s)
	registerDeploymentAPIPaths(s)
	registerVendorAPIPaths(s)
	registerNotificationAPIPaths(s)
	registerResourceAPIPaths(s)
	registerAutomationAPIPaths(s)
	registerFeatureAPIPaths(s)
	registerMetadataAPIPaths(s)
	registerScheduledTasksAPIPaths(s)

	registerAPIv2Paths(s)
}

func registerAuditTrailAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiAuditTrailPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allAuditTrailEvents).Methods("GET")
	s.HandleFunc(core.ApiGetEndpoint, getAuditTrailEntry).Methods("GET")
}

func registerInviteAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiInvitePrefix).Subrouter()
	s.HandleFunc(core.ApiSendInviteEndpoint, sendInviteToOrganization).Methods("POST")
}

func registerVerificationAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiVerificationPrefix).Subrouter()
	s.HandleFunc(core.ApiResendVerificationEndpoint, requestResendUserVerificationEmail).Methods("POST")
}

func registerUserAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiUserPrefix).Subrouter()
	s.Use(webcore.VerifyAPIKeyHasAccessToUser)
	s.HandleFunc(core.ApiUserProfileEndpoint, updateUserProfile).Methods("POST")
	s.HandleFunc(core.ApiUserOrgsEndpoint, getAllOrganizationsForUser).Methods("GET")
}

func registerOrgAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiOrgPrefix).Subrouter()
	s.HandleFunc(core.ApiOrgUsersEndpoint, getAllUsersInOrganization).Methods("GET")
}

func registerProcessFlowAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiProcessFlowUrl).Subrouter()
	s.HandleFunc(core.ApiProcessFlowNewUrl, newProcessFlow).Methods("POST")
	s.HandleFunc(core.ApiProcessFlowGetAllUrl, getAllProcessFlows).Methods("GET")
	s.HandleFunc(core.ApiProcessFlowUpdateUrl, updateProcessFlow).Methods("POST")
	s.HandleFunc(core.ApiProcessFlowGetFullDataUrl, getProcessFlowFullData).Methods("GET")
	s.HandleFunc(core.ApiProcessFlowDeleteEndpoint, deleteProcessFlow).Methods("POST")
}

func registerProcessFlowNodesAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiProcessFlowNodesUrl).Subrouter()
	s.HandleFunc(core.ApiProcessFlowNodesGetTypesUrl, getAllProcessFlowNodeTypes).Methods("GET")
	s.HandleFunc(core.ApiProcessFlowNodesNewUrl, newProcessFlowNode).Methods("POST")
	s.HandleFunc(core.ApiProcessFlowNodesEditUrl, editProcessFlowNode).Methods("POST")
	s.HandleFunc(core.ApiProcessFlowNodesDeleteUrl, deleteProcessFlowNode).Methods("POST")
	s.HandleFunc(core.ApiDuplicateEndpoint, duplicateProcessFlowNode).Methods("POST")

	registerProcessFlowNodeLinksAPIPaths(s)
}

func registerProcessFlowNodeLinksAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiProcessFlowNodeLinksPrefix).Subrouter()
	registerProcessFlowNodeSystemLinksAPIPaths(s)
	registerProcessFlowNodeGLLinksAPIPaths(s)
}

func registerProcessFlowNodeSystemLinksAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITSystemsPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newNodeSystemLink).Methods("POST")
	s.HandleFunc(core.ApiAllEndpoint, allNodeSystemLink).Methods("GET")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteNodeSystemLink).Methods("POST")
}

func registerProcessFlowNodeGLLinksAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiGeneralLedgerPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newNodeGLLink).Methods("POST")
	s.HandleFunc(core.ApiAllEndpoint, allNodeGLLink).Methods("GET")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteNodeGLLink).Methods("POST")
}

func registerProcessFlowIOAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiProcessFlowIOUrl).Subrouter()
	s.HandleFunc(core.ApiProcessFlowIOGetTypesUrl, getAllProcessFlowIOTypes).Methods("GET")
	s.HandleFunc(core.ApiProcessFlowIONewUrl, createNewProcessFlowIO).Methods("POST")
	s.HandleFunc(core.ApiProcessFlowIODeleteUrl, deleteProcessFlowIO).Methods("POST")
	s.HandleFunc(core.ApiProcessFlowIOEditUrl, editProcessFlowIO).Methods("POST")
	s.HandleFunc(core.ApiProcessFlowIOOrderEndpoint, orderProcessFlowIO).Methods("POST")
}

func registerProcessFlowEdgesAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiProcessFlowEdgesUrl).Subrouter()
	s.HandleFunc(core.ApiProcessFlowEdgesNewUrl, createNewProcessFlowEdge).Methods("POST")
	s.HandleFunc(core.ApiProcessFlowEdgesDeleteUrl, deleteProcessFlowEdge).Methods("POST")
}

func registerRiskAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiRiskPrefix).Subrouter()
	s.HandleFunc(core.ApiNewRiskEndpoint, createNewRisk).Methods("POST")
	s.HandleFunc(core.ApiDeleteRiskEndpoint, deleteRisks).Methods("POST")
	s.HandleFunc(core.ApiEditRiskEndpoint, editRisk).Methods("POST")
	s.HandleFunc(core.ApiAddRiskToNodeEndpoint, addRisksToNode).Methods("POST")
	s.HandleFunc(core.ApiGetAllRisksEndpoint, getAllRisks).Methods("GET")
	s.HandleFunc(core.ApiGetSingleRiskEndpoint, getSingleRisk).Methods("GET")

	registerRiskLinkAPIPaths(s)
}

func registerRiskLinkAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiProcessFlowNodeLinksPrefix).Subrouter()
	registerRiskLinkSystemAPIPaths(s)
	registerRiskLinkGLAPIPaths(s)
}

func registerRiskLinkSystemAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITSystemsPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allRiskSystemLinks).Methods("GET")
}

func registerRiskLinkGLAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiGeneralLedgerPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allRiskGeneralLedgerAccountLinks).Methods("GET")
}

func registerControlAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiControlPrefix).Subrouter()
	s.HandleFunc(core.ApiNewControlEndpoint, newControl).Methods("POST")
	s.HandleFunc(core.ApiGetControlTypesEndpoint, getControlTypes).Methods("GET")
	s.HandleFunc(core.ApiDeleteControlEndpoint, deleteControls).Methods("POST")
	s.HandleFunc(core.ApiAddControlEndpoint, addControls).Methods("POST")
	s.HandleFunc(core.ApiEditControlEndpoint, editControl).Methods("POST")
	s.HandleFunc(core.ApiGetAllControlEndpoint, getAllControls).Methods("GET")
	s.HandleFunc(core.ApiGetSingleControlEndpoint, getSingleControl).Methods("GET")

	registerControlLinkAPIPaths(s)
}

func registerControlLinkAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiProcessFlowNodeLinksPrefix).Subrouter()
	registerControlLinkSystemAPIPaths(s)
	registerControlLinkGLAPIPaths(s)
	registerControlLinkDocAPIPaths(s)
}

func registerControlLinkSystemAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITSystemsPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allControlSystemLinks).Methods("GET")
}

func registerControlLinkGLAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiGeneralLedgerPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allControlGeneralLedgerAccountLinks).Methods("GET")
}

func registerControlLinkDocAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiControlDocumentationPrefix).Subrouter()
	registerControlLinkDocCatAPIPaths(s)
	registerControlLinkFolderAPIPaths(s)
}

func registerControlLinkDocCatAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiDocCatPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allControlDocCatLinks).Methods("GET")
}

func registerControlLinkFolderAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiFolderPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allControlFolderLinks).Methods("GET")
}

func registerControlDocumentationAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiControlDocumentationPrefix).Subrouter()
	s.HandleFunc(core.ApiGetControlDocumentationCategoryEndpoint, getControlDocumentationCategory).Methods("GET")
	s.HandleFunc(core.ApiAllControlDocumentationCategoryEndpoint, allControlDocumentationCategories).Methods("GET")
	s.HandleFunc(core.ApiNewControlDocumentationCategoryEndpoint, newControlDocumentationCategory).Methods("POST")
	s.HandleFunc(core.ApiEditControlDocumentationCategoryEndpoint, editControlDocumentationCategory).Methods("POST")
	s.HandleFunc(core.ApiDeleteControlDocumentationCategoryEndpoint, deleteControlDocumentationCategory).Methods("POST")

	s.HandleFunc(core.ApiUploadControlDocumentationEndpoint, uploadControlDocumentation).Methods("POST")
	s.HandleFunc(core.ApiAllControlDocumentationEndpoint, allControlDocumentation).Methods("GET")
	s.HandleFunc(core.ApiDeleteControlDocumentationEndpoint, deleteControlDocumentation).Methods("POST")
	s.HandleFunc(core.ApiDownloadControlDocumentationEndpoint, downloadControlDocumentation).Methods("GET")
	s.HandleFunc(core.ApiGetControlDocumentationEndpoint, getControlDocumentation).Methods("GET")
	s.HandleFunc(core.ApiEditControlDocumentationEndpoint, editControlDocumentation).Methods("POST")
	s.HandleFunc(core.ApiRegenPreviewControlDocumentationEndpoint, regeneratePreview).Methods("POST")

	registerControlDocVersionsAPIPaths(s)
	registerFileFolderAPIPaths(s)
}

func registerControlDocVersionsAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiFileVersionPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allFileVersions).Methods("GET")
	s.HandleFunc(core.ApiGetEndpoint, getFileVersion).Methods("GET")
}

func registerFileFolderAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiFolderPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newFolder).Methods("POST")
	s.HandleFunc(core.ApiUpdateEndpoint, updateFolder).Methods("POST")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteFolder).Methods("POST")

	registerFileFolderLinksAPIPaths(s)
}

func registerFileFolderLinksAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiLinkPrefix).Subrouter()

	registerFileFolderFileLinksAPIPaths(s)
}

func registerFileFolderFileLinksAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiDocFilePrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allFolderFileLinks).Methods("GET")
	s.HandleFunc(core.ApiNewEndpoint, newFolderFileLinks).Methods("POST")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteFolderFileLink).Methods("POST")
}

func registerRoleAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiRolePrefix).Subrouter()
	s.HandleFunc(core.ApiGetOrganizationRolesEndpoint, getAllOrganizationRoles).Methods("GET")
	s.HandleFunc(core.ApiGetSingleRoleEndpoint, getSingleRole).Methods("GET")
	s.HandleFunc(core.ApiNewRoleEndpoint, newRole).Methods("POST")
	s.HandleFunc(core.ApiEditRoleEndpoint, editRole).Methods("POST")
	s.HandleFunc(core.ApiDeleteRoleEndpoint, deleteRole).Methods("POST")
	s.HandleFunc(core.ApiAddUsersToRoleEndpoint, addUsersToRole).Methods("POST")
}

func registerGeneralLedgerAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiGeneralLedgerPrefix).Subrouter()
	s.HandleFunc(core.ApiGetGLLevelEndpoint, getGL).Methods("GET")
	s.HandleFunc(core.ApiNewGLCategoryEndpoint, createNewGLCategory).Methods("POST")
	s.HandleFunc(core.ApiEditGLCategoryEndpoint, editGLCategory).Methods("POST")
	s.HandleFunc(core.ApiDeleteGLCategoryEndpoint, deleteGLCategory).Methods("POST")
	s.HandleFunc(core.ApiNewGLAccountEndpoint, createNewGLAccount).Methods("POST")
	s.HandleFunc(core.ApiGetGLAccountEndpoint, getGLAccount).Methods("GET")
	s.HandleFunc(core.ApiEditGLAccountEndpoint, editGLAccount).Methods("POST")
	s.HandleFunc(core.ApiDeleteGLAccountEndpoint, deleteGLAccount).Methods("POST")
}

func registerITAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITPrefix).Subrouter()
	registerITSystemsAPIPaths(s)
	registerITDbAPIPaths(s)
	registerITServerAPIPaths(s)

	s.HandleFunc(core.ApiITDeleteDbSysLinkEndpoint, deleteDatabaseSystemLink).Methods("POST")
}

func registerITSystemsAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITSystemsPrefix).Subrouter()
	s.HandleFunc(core.ApiITSystemsNewEndpoint, newSystem).Methods("POST")
	s.HandleFunc(core.ApiITSystemsAllEndpoint, getAllSystems).Methods("GET")
	s.HandleFunc(core.ApiITSystemGetEndpoint, getSystem).Methods("GET")
	s.HandleFunc(core.ApiITSystemEditEndpoint, editSystem).Methods("POST")
	s.HandleFunc(core.ApiITSystemDeleteEndpoint, deleteSystem).Methods("POST")
	s.HandleFunc(core.ApiITSystemLinkDbEndpoint, linkDatabasesToSystem).Methods("POST")
}

func registerITDbAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITDbPrefix).Subrouter()
	s.HandleFunc(core.ApiITDbNewEndpoint, newDb).Methods("POST")
	s.HandleFunc(core.ApiITDbAllEndpoint, getAllDb).Methods("GET")
	s.HandleFunc(core.ApiITDbTypesEndpoint, getDbTypes).Methods("GET")
	s.HandleFunc(core.ApiITDbGetEndpoint, getDb).Methods("GET")
	s.HandleFunc(core.ApiITDbEditEndpoint, editDb).Methods("POST")
	s.HandleFunc(core.ApiITDbDeleteEndpoint, deleteDb).Methods("POST")
	s.HandleFunc(core.ApiITDbLinkSysEndpoint, linkSystemsToDatabase).Methods("POST")

	registerITDbConnAPIPaths(s)
	registerITSqlAPIPaths(s)
}

func registerITDbConnAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITDbConnPrefix).Subrouter()
	s.HandleFunc(core.ApiITDbConnNewEndpoint, newDbConnection).Methods("POST")
	s.HandleFunc(core.ApiITDbConnDeleteEndpoint, deleteDbConnection).Methods("POST")
}

func registerITSqlAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITSqlPrefix).Subrouter()

	registerITSqlRefreshAPIPaths(s)
	registerITSqlSchemaAPIPaths(s)
	registerITSqlQueriesAPIPaths(s)
	registerITSqlRequestsAPIPaths(s)
}

func registerITSqlRefreshAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITSqlRefreshPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allDatabaseRefresh).Methods("GET")
	s.HandleFunc(core.ApiGetEndpoint, getDatabaseRefresh).Methods("GET")
	s.HandleFunc(core.ApiNewEndpoint, newDatabaseRefresh).Methods("POST")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteDatabaseRefresh).Methods("POST")
}

func registerITSqlSchemaAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITSqlSchemaPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allDatabaseSchemas).Methods("GET")
	s.HandleFunc(core.ApiGetEndpoint, getDatabaseSchema).Methods("GET")
}

func registerITSqlQueriesAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITSqlQueriesPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allDatabaseQuery).Methods("GET")
	s.HandleFunc(core.ApiGetEndpoint, getDatabaseQuery).Methods("GET")
	s.HandleFunc(core.ApiNewEndpoint, newDatabaseQuery).Methods("POST")
	s.HandleFunc(core.ApiUpdateEndpoint, updateDatabaseQuery).Methods("POST")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteDatabaseQuery).Methods("POST")
	s.HandleFunc(core.ApiRunEndpoint, runDatabaseQuery).Methods("POST")
}

func registerITSqlRequestsAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITSqlRequestsPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newSqlRequest).Methods("POST")
	s.HandleFunc(core.ApiAllEndpoint, allSqlRequest).Methods("GET")
	s.HandleFunc(core.ApiITSqlRequestsStatusEndpoint, statusSqlRequest).Methods("GET")
	s.HandleFunc(core.ApiITSqlRequestsStatusEndpoint, modifyStatusSqlRequest).Methods("POST")
	s.HandleFunc(core.ApiGetEndpoint, getSqlRequest).Methods("GET")
	s.HandleFunc(core.ApiUpdateEndpoint, updateSqlRequest).Methods("POST")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteSqlRequest).Methods("POST")
}

func registerITServerAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITServerPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newServer).Methods("POST")
	s.HandleFunc(core.ApiAllEndpoint, allServers).Methods("GET")
	s.HandleFunc(core.ApiGetEndpoint, getServer).Methods("GET")
	s.HandleFunc(core.ApiUpdateEndpoint, updateServer).Methods("POST")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteServer).Methods("POST")
}

func registerDocRequestsAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiDocRequestPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newDocumentRequest).Methods("POST")
	s.HandleFunc(core.ApiGetEndpoint, getDocumentRequest).Methods("GET")
	s.HandleFunc(core.ApiAllEndpoint, allDocumentRequests).Methods("GET")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteDocumentRequest).Methods("POST")
	s.HandleFunc(core.ApiUpdateEndpoint, updateDocumentRequest).Methods("POST")
	s.HandleFunc(core.ApiDocRequestCompleteEndpoint, completeDocumentRequest).Methods("POST")

	registerDocRequestLinksAPIPaths(s)
}

func registerDocRequestLinksAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiLinkPrefix).Subrouter()
	registerDocRequestDocCatLinksAPIPaths(s)
	registerDocRequestControlLinksAPIPaths(s)
}

func registerDocRequestControlLinksAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiControlPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allDocRequestControlLinks).Methods("GET")
}

func registerDocRequestDocCatLinksAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiDocCatPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allDocRequestDocCatLinks).Methods("GET")
}

func registerCommentsAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiCommentsPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newComment).Methods("POST")
	s.HandleFunc(core.ApiAllEndpoint, allComments).Methods("GET")
	s.HandleFunc(core.ApiUpdateEndpoint, updateComment).Methods("POST")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteComment).Methods("POST")
}

func registerDeploymentAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiDeploymentPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newDeployment).Methods("POST")
	s.HandleFunc(core.ApiUpdateEndpoint, updateDeployment).Methods("POST")

	registerDeploymentLinkAPIPaths(s)
}

func registerDeploymentLinkAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiLinkPrefix).Subrouter()
	registerDeploymentServerLinkAPIPaths(s)
}

func registerDeploymentServerLinkAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITServerPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newDeploymentServerLink).Methods("POST")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteDeploymentServerLink).Methods("POST")
}

func registerVendorAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiVendorPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newVendor).Methods("POST")
	s.HandleFunc(core.ApiAllEndpoint, allVendors).Methods("GET")
	s.HandleFunc(core.ApiGetEndpoint, getVendor).Methods("GET")
	s.HandleFunc(core.ApiUpdateEndpoint, updateVendor).Methods("POST")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteVendor).Methods("POST")

	registerVendorProductAPIPaths(s)
}

func registerVendorProductAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiVendorProductPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newVendorProduct).Methods("POST")
	s.HandleFunc(core.ApiAllEndpoint, allVendorProducts).Methods("GET")
	s.HandleFunc(core.ApiGetEndpoint, getVendorProduct).Methods("GET")
	s.HandleFunc(core.ApiUpdateEndpoint, updateVendorProduct).Methods("POST")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteVendorProduct).Methods("POST")

	registerVendorProductSocAPIPaths(s)
}

func registerVendorProductSocAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiVendorProductSocPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, linkVendorProductSoc).Methods("POST")
	s.HandleFunc(core.ApiDeleteEndpoint, unlinkVendorProductSoc).Methods("POST")
}

func registerNotificationAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiNotificationPrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allNotifications).Methods("GET")
	s.HandleFunc(core.ApiReadEndpoint, markNotificationRead).Methods("POST")
}

func registerResourceAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiResourcePrefix).Subrouter()
	s.HandleFunc(core.ApiGetEndpoint, getResource).Methods("GET")
}

func registerAutomationAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.DashboardAutomationPrefix).Subrouter()
	registerDataAPIPaths(s)
	registerCodeAPIPaths(s)
	registerScriptsAPIPaths(s)
	registerLogsAPIPaths(s)
}

func registerDataAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.DashboardDataPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newClientData).Methods("POST")
	s.HandleFunc(core.ApiUpdateEndpoint, updateClientData).Methods("POST")
	s.HandleFunc(core.ApiAllEndpoint, allClientData).Methods("GET")
	s.HandleFunc(core.ApiGetEndpoint, getClientData).Methods("GET")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteClientData).Methods("POST")

	registerDataSourceAPIPaths(s)
}

func registerDataSourceAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.DashboardDataSourcePrefix).Subrouter()
	s.HandleFunc(core.ApiAllEndpoint, allDataSourceOptions).Methods("GET")
	s.HandleFunc(core.ApiGetEndpoint, getDataSource).Methods("GET")
}

func registerCodeAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.DashboardCodePrefix).Subrouter()
	s.HandleFunc(core.ApiSaveEndpoint, saveCode).Methods("POST")
	s.HandleFunc(core.ApiGetEndpoint, getCode).Methods("GET")
	s.HandleFunc(core.ApiAllEndpoint, allCode).Methods("GET")
	s.HandleFunc(core.ApiLinkEndpoint, getCodeLink).Methods("GET")

	registerCodeStatusAPIPaths(s)
	registerCodeRunAPIPaths(s)
}

func registerCodeRunAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.DashboardRunPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, runCode).Methods("POST").Name(webcore.ApiRunCodeRouteName)
	s.HandleFunc(core.ApiAllEndpoint, allCodeRuns).Methods("GET")
	s.HandleFunc(core.ApiGetEndpoint, getCodeRun).Methods("GET")
}

func registerCodeStatusAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.DashboardStatusPrefix).Subrouter()
	s.HandleFunc(core.ApiGetEndpoint, getCodeBuildStatus).Methods("GET")
}

func registerScriptsAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.DashboardScriptPrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, newClientScript).Methods("POST")
	s.HandleFunc(core.ApiAllEndpoint, allClientScripts).Methods("GET")
	s.HandleFunc(core.ApiGetEndpoint, getClientScript).Methods("GET")
	s.HandleFunc(core.ApiDeleteEndpoint, deleteClientScript).Methods("POST")
	s.HandleFunc(core.ApiUpdateEndpoint, updateClientScript).Methods("POST")

	registerScriptLinksAPIPaths(s)
}

func registerScriptLinksAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.DashboardLinkPrefix).Subrouter()
	s.HandleFunc(core.ApiGetEndpoint, getScriptCodeLink).Methods("GET")
}

func registerLogsAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.DashboardLogsPrefix).Subrouter()
	s.HandleFunc(core.ApiGetEndpoint, getLogs).Methods("GET")
}

func registerFeatureAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiFeaturePrefix).Subrouter()
	s.HandleFunc(core.ApiNewEndpoint, enableFeature).Methods("POST")
}

func registerMetadataAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiMetadataPrefix).Subrouter()
	s.HandleFunc("/paramTypes", getParamTypesMetadata).Methods("GET")
}

func registerScheduledTasksAPIPaths(r *mux.Router) {
	s := r.PathPrefix("/schedule").Subrouter()
	s.HandleFunc("/", allScheduledTasks).Methods("GET")
}
