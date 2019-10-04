package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/render"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
	"strings"
)

func postLogin(w http.ResponseWriter, r *http.Request) {
	var err error

	// Expect an "email" field and "nonce" field in the post data encoded as x-www-form-urlencoded.
	// The "nonce" field is a CRSF token that will be verified twice. Once here and once after
	// Okta redirects back to us.
	//
	// In this function, we take the email, strip out its domain, see if it matches any IdP we know
	// stored in our database.
	// If it does match, then tell the client to redirect to the appropriate login page.
	//
	// If it doesn't match, then either tell the user to "Get Started" or in certain
	// situations register. Note that chances are we only want to register external auditors
	// with whom we don't have SAML support for.
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	// Retrieve the email address and parse the domain.
	if err = r.ParseForm(); err != nil || len(r.PostForm) == 0 {
		core.Warning("Failed to parse form data: " + core.ErrorString(err))
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	emailData := r.PostForm["email"]
	csrfToken := r.PostForm["csrf"]

	if len(emailData) == 0 || len(csrfToken) == 0 {
		core.Warning("Empty email or CSRF.")
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	if ok, err := webcore.VerifyCSRFToken(csrfToken[0], r); !ok || err != nil {
		core.Warning("Failed CSRF check: " + core.ErrorString(err))
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	// I don't think performing an email validation here is necessary since all we really
	// want to do is to parse out the domain. If it matches something it does and if it doesn't
	// it doesn't. No need to detect if the user entered a valid email.
	email := strings.TrimSpace(emailData[0])
	_, domain := core.ParseEmailAddress(email)

	// Find the IdP. If not found, return an error. Use Error 400 but include a boolean
	// flag in the body to indicate that we don't know what domain this is. If it is found
	// return OK 200 and include the desirable login URL in the body.
	var idpIden string
	idpIden, err = database.FindSAMLIdPFromDomain(domain)

	// This error should only crop up if something went terribly wrong on our side.
	if err != nil {
		core.Warning("Find SAML IdP Error: " + err.Error() + " (" + domain + ")")
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}

	if idpIden == "" {
		core.Warning("Can not found IdP: " + domain)
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct {
			CanNotFindIdP bool
		}{
			true,
		})
		return
	}

	core.Info(webcore.CreateOktaLoginUrl(idpIden, csrfToken[0], "filler"))
	jsonWriter.Encode(struct {
		LoginUrl string
	}{
		// Pass the CSRF token as the state and verify it upon redirect
		// because why not.
		webcore.CreateOktaLoginUrl(idpIden, csrfToken[0], "filler"),
	})
}

func getSamlLoginCallbackError(prefix string, err error, w http.ResponseWriter, r *http.Request) {
	core.Warning(prefix + " :: " + core.ErrorString(err))
	webcore.ClearCSRFTokenFromSession(w, r)
	render.RetrieveTemplate(render.RedirectTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			render.CreateRedirectParams(w, r, "Oops!",
				"Something went wrong! Please try again.",
				webcore.MustGetRouteUrl(webcore.LoginRouteName)))
}

func getSamlLoginCallback(w http.ResponseWriter, r *http.Request) {
	// There are two query params: 'state' and 'code'.
	// 'state' will contain the CSRF token. We should check this against the user's cookie.
	// 'code' will contain the authorization code which we need to convert into a token.
	query := r.URL.Query()
	state, sok := query["state"]
	code, cok := query["code"]

	if !sok || !cok || len(state) == 0 || len(code) == 0 {
		getSamlLoginCallbackError("Empty state or code.", nil, w, r)
		return
	}

	if ok, err := webcore.VerifyCSRFToken(state[0], r); !ok || err != nil {
		getSamlLoginCallbackError("Failed CSRF check for SAML Login Callback.", err, w, r)
		return
	}
	webcore.ClearCSRFTokenFromSession(w, r)

	// Retrieve the access/ID token from Okta and redirect if successful.
	// Note that core.OktaObtainTokens will store the tokens where necessary.
	tokens, err := webcore.OktaObtainTokens(code[0], false)
	if err != nil {
		getSamlLoginCallbackError("Failed to obtain OIDC tokens.", err, w, r)
		return
	}

	session, err := webcore.CreateUserSessionFromTokens(tokens, r)
	if err != nil {
		getSamlLoginCallbackError("Failed to create user session.", err, w, r)
		return
	}

	// Store session in the database.
	if err = database.StoreUserSession(session); err != nil {
		getSamlLoginCallbackError("Failed to store user session (server).", err, w, r)
		return
	}

	// Store session id as a cookie.
	if err = webcore.StoreUserSessionOnClient(session, w); err != nil {
		getSamlLoginCallbackError("Failed to store user session (client).", err, w, r)
		return
	}

	http.Redirect(w, r, webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName), http.StatusFound)
}

func getLogoutCallbackError(prefix string, err error, w http.ResponseWriter, r *http.Request) {
	core.Warning(prefix + " :: " + core.ErrorString(err))
	webcore.ClearCSRFTokenFromSession(w, r)
	render.RetrieveTemplate(render.RedirectTemplateKey).
		ExecuteTemplate(
			w,
			"base",
			render.CreateRedirectParams(w, r, "Oops!",
				"Something went wrong! Please try again.",
				webcore.MustGetRouteUrl(webcore.DashboardHomeRouteName)))
}

func getLogout(w http.ResponseWriter, r *http.Request) {
	queryVals := r.URL.Query()
	csrfToken, ok := queryVals["csrf"]
	if !ok || len(csrfToken) == 0 {
		getLogoutCallbackError("Failed to logout (no csrf)", nil, w, r)
		return
	}

	if ok, err := webcore.VerifyCSRFToken(csrfToken[0], r); !ok || err != nil {
		getLogoutCallbackError("Failed to logout (bad csrf)", err, w, r)
		return
	}
	webcore.ClearCSRFTokenFromSession(w, r)

	session, err := webcore.FindSessionInContext(r.Context())
	if err != nil {
		getLogoutCallbackError("Failed to logout (no sess)", err, w, r)
		return
	}

	// Need to do a few things on logout:
	// 	1) Delete user session cookie.
	// 	2) Delete user session in database
	// 	3) Delete user session with Okta.
	webcore.DeleteUserSessionOnClient(w)

	if err = database.DeleteUserSession(session.SessionId); err != nil {
		getLogoutCallbackError("Failed to logout (server)", err, w, r)
		return
	}

	http.Redirect(w, r, webcore.CreateOktaLogoutUrl(session.IdToken), http.StatusTemporaryRedirect)
}
