package rest

import (
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
)

func RegisterPaths(r *mux.Router) {
	r.HandleFunc(core.GetStartedUrl, postGettingStartedInterest).Methods("POST").Name(string(webcore.GettingStartedPostRouteName))

	// core.LoginURL is the POST request the user will send with just their email.
	// core.SamlCallbackUrl is the GET request the user's SAML IdP will redirect to upon
	// successful login.
	r.HandleFunc(core.LoginUrl, postLogin).Methods("POST").Name(string(webcore.LoginPostRouteName))
	r.HandleFunc(core.SamlCallbackUrl, getSamlLoginCallback).Methods("GET").Name(string(webcore.SamlCallbackRouteName))
	r.HandleFunc(core.LogoutUrl, getLogout).Methods("GET").Name(string(webcore.LogoutRouteName))
}
