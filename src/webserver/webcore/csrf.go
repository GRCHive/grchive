package webcore

import (
	"github.com/google/uuid"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"net/http"
)

// Add the CSRF token to the session (cookie) and also adds it to the input map[string]interface{}
// so that the templating engine/frontend can add it to the HTML/JS as needed for verification.
func AddCSRFTokenToRequest(w http.ResponseWriter, r *http.Request, pageVars map[string]interface{}) (map[string]interface{}, error) {
	newPageVars := core.CopyMap(pageVars)
	session, err := ClientShortSessionStore.Get(r, "csrf")
	if err != nil {
		core.Warning("Failed to retrieve from session: " + err.Error())
		return pageVars, err
	}

	// This is probably ok to use instead of uuid.NewRandom. I don't think we'll
	// encounter where uuid.Must will fail? Famous last words.
	token := uuid.New().String()
	session.Values["csrf"] = token
	newPageVars["Csrf"] = token
	err = session.Save(r, w)
	if err != nil {
		core.Warning("Failed to save to session: " + err.Error())
		return pageVars, err
	}
	return newPageVars, nil
}

func GetCSRFToken(r *http.Request) (string, error) {
	csrfToken := r.FormValue("csrf")
	return csrfToken, nil
}

func VerifyCSRFToken(token string, r *http.Request) (bool, error) {
	session, err := ClientShortSessionStore.Get(r, "csrf")
	if err != nil {
		return false, err
	}
	return session.Values["csrf"] == token, nil
}

func ClearCSRFTokenFromSession(w http.ResponseWriter, r *http.Request) {
	session, err := ClientShortSessionStore.Get(r, "csrf")
	if err != nil {
		core.Warning("Failed to retrieve from session: " + err.Error())
		return
	}
	session.Values["csrf"] = ""
	err = session.Save(r, w)
	if err != nil {
		core.Warning("Failed to save to session: " + err.Error())
		return
	}
}
