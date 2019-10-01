package rest

import (
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func RegisterPaths(r *mux.Router) {
	r.HandleFunc(core.GetStartedUrl, postGettingStartedInterest).Methods("POST")

	// core.LoginURL is the POST request the user will send with just their email.
	// core.SamlCallbackUrl is the GET request the user's SAML IdP will redirect to upon
	// successful login.
	r.HandleFunc(core.LoginUrl, postLogin).Methods("POST")
	r.HandleFunc(core.SamlCallbackUrl, getSamlLoginCallback).Methods("GET")
}
