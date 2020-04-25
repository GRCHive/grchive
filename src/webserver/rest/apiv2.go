package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/webcore"
)

func registerAPIv2Paths(r *mux.Router) {
	s := r.PathPrefix("/v2").Subrouter()

	registerAPIv2RequestsPaths(s)
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
