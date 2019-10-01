package core

import (
	"github.com/google/uuid"
	"net/http"
)

// Add the CSRF token to the session (cookie) and also adds it to the input map[string]interface{}
// so that the templating engine/frontend can add it to the HTML/JS as needed for verification.
func AddCSRFTokenToRequest(w http.ResponseWriter, r *http.Request, pageVars map[string]interface{}) (map[string]interface{}, error) {
	newPageVars := CopyMap(pageVars)
	session, err := SessionStore.Get(r, "csrf")
	if err != nil {
		Warning("Failed to retrieve from session: " + err.Error())
		return pageVars, err
	}

	// This is probably ok to use instead of uuid.NewRandom. I don't think we'll
	// encounter where uuid.Must will fail? Famous last words.
	token := uuid.New().String()
	session.Values["csrf"] = token
	newPageVars["Csrf"] = token
	err = session.Save(r, w)
	if err != nil {
		Warning("Failed to save to session: " + err.Error())
		return pageVars, err
	}
	return newPageVars, nil
}

func VerifyCSRFToken(token string, r *http.Request) (bool, error) {
	session, err := SessionStore.Get(r, "csrf")
	if err != nil {
		return false, err
	}

	return session.Values["csrf"] == token, nil
}
