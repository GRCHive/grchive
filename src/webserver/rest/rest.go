package rest

import (
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
)

func RegisterPaths(r *mux.Router) {
	altApiRouter := r.NewRoute().Subrouter()
	altApiRouter.Use(webcore.ObtainUserSessionInContextMiddleware)

	altApiRouter.HandleFunc(core.GetStartedUrl, postGettingStartedInterest).Methods("POST").Name(string(webcore.GettingStartedPostRouteName))

	// core.LoginURL is the POST request the user will send with just their email.
	// core.CreateSamlCallbackUrl() is the GET request the user's SAML IdP will redirect to upon
	// successful login.
	altApiRouter.HandleFunc(core.LoginUrl, postLogin).Methods("POST").Name(string(webcore.LoginPostRouteName))
	altApiRouter.HandleFunc(core.RegisterUrl, postRegister).Methods("POST").Name(string(webcore.RegisterPostRouteName))
	altApiRouter.HandleFunc(core.CreateSamlCallbackUrl(), getSamlLoginCallback).Methods("GET").Name(string(webcore.SamlCallbackRouteName))
	altApiRouter.HandleFunc(core.LogoutUrl, getLogout).Methods("GET").Name(string(webcore.LogoutRouteName))
	altApiRouter.HandleFunc(core.VerifyEmailUrl, verifyUserEmail).Methods("GET").Name(webcore.EmailVerifyRouteName)
	altApiRouter.HandleFunc(core.AcceptInviteUrl, acceptInviteToOrganization).Methods("GET").Name(webcore.AcceptInviteRouteName)

	// REST API
	registerAPIPaths(r)
}

func registerAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiUrl).Subrouter()
	s.Use(webcore.ObtainAPIKeyRoleInContextMiddleware)

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
}

func registerInviteAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiInvitePrefix).Subrouter()
	s.HandleFunc(core.ApiSendInviteEndpoint, sendInviteToOrganization).Methods("POST").Name(webcore.SendInviteRouteName)
}

func registerVerificationAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiVerificationPrefix).Subrouter()
	s.HandleFunc(core.ApiResendVerificationEndpoint, requestResendUserVerificationEmail).Methods("POST").Name(webcore.ResendVerificationRouteName)
}

func registerUserAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiUserPrefix).Subrouter()
	s.HandleFunc(core.ApiUserProfileEndpoint, updateUserProfile).Methods("POST").Name(webcore.UserProfileEditRouteName)
	s.HandleFunc(core.ApiUserOrgsEndpoint, getAllOrganizationsForUser).Methods("GET").Name(webcore.UserGetOrgsRouteName)
}

func registerOrgAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiOrgPrefix).Subrouter()
	s.HandleFunc(core.ApiOrgUsersEndpoint, getAllUsersInOrganization).Methods("GET").Name(webcore.GetAllOrgUsersRouteName)
}

func registerProcessFlowAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiProcessFlowUrl).Subrouter()
	s.HandleFunc(core.ApiProcessFlowNewUrl, newProcessFlow).Methods("POST").Name(webcore.NewProcessFlowRouteName)
	s.HandleFunc(core.ApiProcessFlowGetAllUrl, getAllProcessFlows).Methods("GET").Name(webcore.GetAllProcessFlowRouteName)
	s.HandleFunc(core.ApiProcessFlowUpdateUrl, updateProcessFlow).Methods("POST").Name(webcore.UpdateProcessFlowRouteName)
	s.HandleFunc(core.ApiProcessFlowGetFullDataUrl, getProcessFlowFullData).Methods("GET").Name(webcore.GetProcessFlowFullDataRouteName)
	s.HandleFunc(core.ApiProcessFlowDeleteEndpoint, deleteProcessFlow).Methods("POST").Name(webcore.DeleteProcessFlowRouteName)
}

func registerProcessFlowNodesAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiProcessFlowNodesUrl).Subrouter()
	s.HandleFunc(core.ApiProcessFlowNodesGetTypesUrl, getAllProcessFlowNodeTypes).Methods("GET").Name(webcore.GetAllProcessFlowNodeTypesRouteName)
	s.HandleFunc(core.ApiProcessFlowNodesNewUrl, newProcessFlowNode).Methods("POST").Name(webcore.NewProcessFlowNodeRouteName)
	s.HandleFunc(core.ApiProcessFlowNodesEditUrl, editProcessFlowNode).Methods("POST").Name(webcore.EditProcessFlowNodeRouteName)
	s.HandleFunc(core.ApiProcessFlowNodesDeleteUrl, deleteProcessFlowNode).Methods("POST").Name(webcore.DeleteProcessFlowNodeRouteName)
}

func registerProcessFlowIOAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiProcessFlowIOUrl).Subrouter()
	s.HandleFunc(core.ApiProcessFlowIOGetTypesUrl, getAllProcessFlowIOTypes).Methods("GET").Name(webcore.GetAllProcessFlowIOTypesRouteName)
	s.HandleFunc(core.ApiProcessFlowIONewUrl, createNewProcessFlowIO).Methods("POST").Name(webcore.CreateNewProcessFlowIOTypesRouteName)
	s.HandleFunc(core.ApiProcessFlowIODeleteUrl, deleteProcessFlowIO).Methods("POST").Name(webcore.DeleteProcessFlowIORouteName)
	s.HandleFunc(core.ApiProcessFlowIOEditUrl, editProcessFlowIO).Methods("POST").Name(webcore.EditProcessFlowIORouteName)
}

func registerProcessFlowEdgesAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiProcessFlowEdgesUrl).Subrouter()
	s.HandleFunc(core.ApiProcessFlowEdgesNewUrl, createNewProcessFlowEdge).Methods("POST").Name(webcore.CreateNewProcessFlowEdgeRouteName)
	s.HandleFunc(core.ApiProcessFlowEdgesDeleteUrl, deleteProcessFlowEdge).Methods("POST").Name(webcore.DeleteProcessFlowEdgeRouteName)
}

func registerRiskAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiRiskPrefix).Subrouter()
	s.HandleFunc(core.ApiNewRiskEndpoint, createNewRisk).Methods("POST").Name(webcore.NewRiskRouteName)
	s.HandleFunc(core.ApiDeleteRiskEndpoint, deleteRisks).Methods("POST").Name(webcore.DeleteRiskRouteName)
	s.HandleFunc(core.ApiEditRiskEndpoint, editRisk).Methods("POST").Name(webcore.EditRiskRouteName)
	s.HandleFunc(core.ApiAddRiskToNodeEndpoint, addRisksToNode).Methods("POST").Name(webcore.AddRiskToNodeRouteName)
	s.HandleFunc(core.ApiGetAllRisksEndpoint, getAllRisks).Methods("GET").Name(webcore.GetAllRiskRouteName)
	s.HandleFunc(core.ApiGetSingleRiskEndpoint, getSingleRisk).Methods("GET").Name(webcore.GetSingleRiskRouteName)
}

func registerControlAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiControlPrefix).Subrouter()
	s.HandleFunc(core.ApiNewControlEndpoint, newControl).Methods("POST").Name(webcore.NewControlRouteName)
	s.HandleFunc(core.ApiGetControlTypesEndpoint, getControlTypes).Methods("GET").Name(webcore.ControlTypesRouteName)
	s.HandleFunc(core.ApiDeleteControlEndpoint, deleteControls).Methods("POST").Name(webcore.DeleteControlRouteName)
	s.HandleFunc(core.ApiAddControlEndpoint, addControls).Methods("POST").Name(webcore.AddControlRouteName)
	s.HandleFunc(core.ApiEditControlEndpoint, editControl).Methods("POST").Name(webcore.EditControlRouteName)
	s.HandleFunc(core.ApiGetAllControlEndpoint, getAllControls).Methods("GET").Name(webcore.GetAllControlRouteName)
	s.HandleFunc(core.ApiGetSingleControlEndpoint, getSingleControl).Methods("GET").Name(webcore.GetSingleControlRouteName)
}

func registerControlDocumentationAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiControlDocumentationPrefix).Subrouter()
	s.HandleFunc(core.ApiNewControlDocumentationCategoryEndpoint, newControlDocumentationCategory).Methods("POST").Name(webcore.NewControlDocCatRouteName)
	s.HandleFunc(core.ApiEditControlDocumentationCategoryEndpoint, editControlDocumentationCategory).Methods("POST").Name(webcore.EditControlDocCatRouteName)
	s.HandleFunc(core.ApiDeleteControlDocumentationCategoryEndpoint, deleteControlDocumentationCategory).Methods("POST").Name(webcore.DeleteControlDocCatRouteName)
	s.HandleFunc(core.ApiUploadControlDocumentationEndpoint, uploadControlDocumentation).Methods("POST").Name(webcore.UploadControlDocRouteName)
	s.HandleFunc(core.ApiGetControlDocumentationEndpoint, getControlDocumentation).Methods("GET").Name(webcore.UploadControlDocRouteName)
	s.HandleFunc(core.ApiDeleteControlDocumentationEndpoint, deleteControlDocumentation).Methods("POST").Name(webcore.UploadControlDocRouteName)
	s.HandleFunc(core.ApiDownloadControlDocumentationEndpoint, downloadControlDocumentation).Methods("GET").Name(webcore.UploadControlDocRouteName)
}

func registerRoleAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiRolePrefix).Subrouter()
	s.HandleFunc(core.ApiGetOrganizationRolesEndpoint, getAllOrganizationRoles).Methods("GET").Name(webcore.GetOrgRolesRouteName)
	s.HandleFunc(core.ApiGetSingleRoleEndpoint, getSingleRole).Methods("GET").Name(webcore.GetSingleRoleRouteName)
	s.HandleFunc(core.ApiNewRoleEndpoint, newRole).Methods("POST").Name(webcore.NewRoleRouteName)
	s.HandleFunc(core.ApiEditRoleEndpoint, editRole).Methods("POST").Name(webcore.EditRoleRouteName)
	s.HandleFunc(core.ApiDeleteRoleEndpoint, deleteRole).Methods("POST").Name(webcore.DeleteRoleRouteName)
	s.HandleFunc(core.ApiAddUsersToRoleEndpoint, addUsersToRole).Methods("POST").Name(webcore.AddUsersToRoleRouteName)
}

func registerGeneralLedgerAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiGeneralLedgerPrefix).Subrouter()
	s.HandleFunc(core.ApiGetGLLevelEndpoint, getGL).Methods("GET").Name(webcore.ApiGetGLRouteName)
	s.HandleFunc(core.ApiNewGLCategoryEndpoint, createNewGLCategory).Methods("POST").Name(webcore.ApiNewGLCatRouteName)
	s.HandleFunc(core.ApiEditGLCategoryEndpoint, editGLCategory).Methods("POST").Name(webcore.ApiEditGLCatRouteName)
	s.HandleFunc(core.ApiDeleteGLCategoryEndpoint, deleteGLCategory).Methods("POST").Name(webcore.ApiDeleteGLCatRouteName)
	s.HandleFunc(core.ApiNewGLAccountEndpoint, createNewGLAccount).Methods("POST").Name(webcore.ApiNewGLAccRouteName)
	s.HandleFunc(core.ApiGetGLAccountEndpoint, getGLAccount).Methods("GET").Name(webcore.ApiGetGLAccRouteName)
	s.HandleFunc(core.ApiEditGLAccountEndpoint, editGLAccount).Methods("POST").Name(webcore.ApiEditGLAccRouteName)
	s.HandleFunc(core.ApiDeleteGLAccountEndpoint, deleteGLAccount).Methods("POST").Name(webcore.ApiDeleteGLAccRouteName)
}

func registerITAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITPrefix).Subrouter()
	registerITSystemsAPIPaths(s)
	registerITDbAPIPaths(s)
	registerITInfraAPIPaths(s)
}

func registerITSystemsAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITSystemsPrefix).Subrouter()
	s.HandleFunc(core.ApiITSystemsNewEndpoint, newSystem).Methods("POST").Name(webcore.ApiNewSystemRouteName)
	s.HandleFunc(core.ApiITSystemsAllEndpoint, getAllSystems).Methods("GET").Name(webcore.ApiSystemAllRouteName)
	s.HandleFunc(core.ApiITSystemGetEndpoint, getSystem).Methods("GET").Name(webcore.ApiGetSystemRouteName)
	s.HandleFunc(core.ApiITSystemEditEndpoint, editSystem).Methods("POST").Name(webcore.ApiEditSystemRouteName)
	s.HandleFunc(core.ApiITSystemDeleteEndpoint, deleteSystem).Methods("POST").Name(webcore.ApiDeleteSystemRouteName)
}

func registerITDbAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiITDbPrefix).Subrouter()
	s.HandleFunc(core.ApiITDbNewEndpoint, newDb).Methods("POST").Name(webcore.ApiNewDbRouteName)
	s.HandleFunc(core.ApiITDbAllEndpoint, getAllDb).Methods("GET").Name(webcore.ApiAllDbRouteName)
	s.HandleFunc(core.ApiITDbTypesEndpoint, getDbTypes).Methods("GET").Name(webcore.ApiTypesDbRouteName)
	s.HandleFunc(core.ApiITDbGetEndpoint, getDb).Methods("GET").Name(webcore.ApiGetDbRouteName)
	s.HandleFunc(core.ApiITDbEditEndpoint, editDb).Methods("POST").Name(webcore.ApiEditDbRouteName)
	s.HandleFunc(core.ApiITDbDeleteEndpoint, deleteDb).Methods("POST").Name(webcore.ApiDeleteDbRouteName)
}

func registerITInfraAPIPaths(r *mux.Router) {
	//s := r.PathPrefix(core.ApiITInfraPrefix).Subrouter()
}
