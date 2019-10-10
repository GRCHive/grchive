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
	registerUserAPIPaths(s)
	registerProcessFlowAPIPaths(s)
	registerProcessFlowNodesAPIPaths(s)
}

func registerUserAPIPaths(r *mux.Router) {
	s := r.PathPrefix(core.ApiUserUrl).Subrouter()
	s.Use(webcore.CreateVerifyUserHasAccessToUserMiddleware(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	s.HandleFunc(core.ApiUserProfileUrl, updateUserProfile).Methods("POST").Name(webcore.UserProfileEditRouteName)
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
}
