package rest

import (
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

func RegisterPaths(r *mux.Router) {
	r.HandleFunc(core.GetStartedUrl, postGettingStartedInterest).Methods("POST").Name(string(webcore.GettingStartedPostRouteName))

	// core.LoginURL is the POST request the user will send with just their email.
	// core.SamlCallbackUrl is the GET request the user's SAML IdP will redirect to upon
	// successful login.
	r.HandleFunc(core.LoginUrl, postLogin).Methods("POST").Name(string(webcore.LoginPostRouteName))
	r.HandleFunc(core.SamlCallbackUrl, getSamlLoginCallback).Methods("GET").Name(string(webcore.SamlCallbackRouteName))
	r.HandleFunc(core.LogoutUrl, getLogout).Methods("GET").Name(string(webcore.LogoutRouteName))

	// REST API
	registerAPIPaths(r)
}

func registerAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiUrl).Subrouter()
	// TODO: API Key verification? For now just verify CSRF.
	// TODO: IP rate limiting?
	s.Use(webcore.CreateVerifyCSRFMiddleware(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	s.Use(webcore.CreateObtainOrganizationInfoFromUserInContextMiddleware(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))

	registerUserAPIPaths(s)
	registerOrgAPIPaths(s)
	registerProcessFlowAPIPaths(s)
	registerProcessFlowNodesAPIPaths(s)
	registerProcessFlowIOAPIPaths(s)
	registerProcessFlowEdgesAPIPaths(s)
	registerRiskAPIPaths(s)
}

func registerUserAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiUserUrl).Subrouter()
	s.Use(webcore.CreateVerifyUserHasAccessToUserMiddleware(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	s.HandleFunc(core.ApiUserProfileUrl, updateUserProfile).Methods("POST").Name(webcore.UserProfileEditRouteName)
}

func registerOrgAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiOrgPrefix).Subrouter()
	s.Use(webcore.CreateVerifyUserHasAccessToOrganizationMiddleware(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	s.HandleFunc(core.ApiOrgUsersEndpoint, getAllUsersInOrganization).Methods("GET").Name(webcore.GetAllOrgUsersRouteName)
}

func registerProcessFlowAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiProcessFlowUrl).Subrouter()
	s.HandleFunc(core.ApiProcessFlowNewUrl, newProcessFlow).Methods("POST").Name(webcore.NewProcessFlowRouteName)
	s.HandleFunc(core.ApiProcessFlowGetAllUrl, getAllProcessFlows).Methods("GET").Name(webcore.GetAllProcessFlowRouteName)
	s.HandleFunc(core.ApiProcessFlowUpdateUrl, updateProcessFlow).Methods("POST").Name(webcore.UpdateProcessFlowRouteName)
	s.HandleFunc(core.ApiProcessFlowGetFullDataUrl, getProcessFlowFullData).Methods("GET").Name(webcore.GetProcessFlowFullDataRouteName)
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
	s.HandleFunc(core.ApiAddRiskToNodeEndpoint, addRisksToNode).Methods("POST").Name(webcore.AddRiskToNodeRouteName)
}
