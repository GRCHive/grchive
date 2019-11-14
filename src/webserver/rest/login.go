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

type BaseLoginInputs struct {
	Email string `webcore:"email"`
	Csrf  string `webcore:"csrf"`
}

type BaseRegisterInputs struct {
	Email      string `webcore:"email"`
	FirstName  string `webcore:"firstName"`
	LastName   string `webcore:"lastName"`
	Password   string `webcore:"password"`
	InviteCode string `webcore:"inviteCode"`
	Csrf       string `webcore:"csrf"`
}

func postRegister(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := BaseRegisterInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Failed to parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ok, err := webcore.VerifyCSRFToken(inputs.Csrf, r); !ok || err != nil {
		core.Warning("Failed CSRF check: " + core.ErrorString(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invite, err := database.FindInviteCodeFromHash(inputs.InviteCode, inputs.Email, core.ServerRole)
	if err != nil {
		core.Warning("Failed invite code verification: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser := core.User{
		FirstName: inputs.FirstName,
		LastName:  inputs.LastName,
		Email:     inputs.Email,
	}

	// Create the user using the Okta API.
	// TODO: We probably need to check what kind of error this is since if
	// 		 it fails because the user was created already it could mean that
	// 		 we failed when creating the user on our end.
	err = webcore.OktaRegisterUserWithPassword(&newUser, inputs.Password)
	if err != nil {
		core.Warning("Failed to register user [okta]: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Store user in database.
	err = webcore.CreateNewUser(&newUser)
	if err != nil {
		core.Warning("Failed to register user [self]: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Mark invite code as being used.
	// It's OK to have this fail since it's not critical for future steps.
	// We should have a cron job or something that tries to fix anything
	// that errors out here.
	err = database.MarkInviteAsUsed(invite)
	if err != nil {
		core.Warning("Failed to mark invite as used: " + err.Error())
	}

	// Redirect user back to login again.
	jsonWriter.Encode(struct {
		RedirectUrl string
	}{
		RedirectUrl: webcore.MustGetRouteUrl(webcore.LoginRouteName),
	})
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	var err error

	// Expect an "email" field and "nonce" field in the post data encoded as x-www-form-urlencoded.
	// The "nonce" field is a CRSF token that will be verified here. TODO: Use middleware?
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

	inputs := BaseLoginInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Failed to parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ok, err := webcore.VerifyCSRFToken(inputs.Csrf, r); !ok || err != nil {
		core.Warning("Failed CSRF check: " + core.ErrorString(err))
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	// I don't think performing an email validation here is necessary since all we really
	// want to do is to parse out the domain. If it matches something it does and if it doesn't
	// it doesn't. No need to detect if the user entered a valid email.
	email := strings.TrimSpace(inputs.Email)
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

	jsonWriter.Encode(struct {
		LoginUrl string
	}{
		// Pass the CSRF token as the state and verify it upon redirect
		// because why not.
		webcore.CreateOktaLoginUrl(idpIden, inputs.Csrf, "filler"),
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
