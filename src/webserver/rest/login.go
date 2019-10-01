package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
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

	if ok, err := core.VerifyCSRFToken(csrfToken[0], r); !ok || err != nil {
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

	jsonWriter.Encode(struct {
		LoginUrl string
	}{
		// Pass the CSRF token as the nonce as well as the state and verify both upon redirect
		// because why not.
		core.CreateOktaLoginUrl(idpIden, csrfToken[0], csrfToken[0]),
	})
}
