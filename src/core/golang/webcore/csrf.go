package webcore

import (
	"github.com/google/uuid"
	"gitlab.com/grchive/grchive/core"
	"net/http"
)

// Add the CSRF token to the secure session (cookie) as well as an unencrypted cookie.
// The unencrypted cookie is for the client to obtain the CSRF token so they can send it
// back to us for the double submit pattern.
func AddCSRFTokenToRequest(w http.ResponseWriter, r *http.Request) error {
	session, err := ClientShortSessionStore.Get(r, "csrf")
	if err != nil {
		core.Warning("Failed to retrieve from session: " + err.Error())
		return err
	}

	// This is probably ok to use instead of uuid.NewRandom. I don't think we'll
	// encounter where uuid.Must will fail? Famous last words.
	token := uuid.New().String()
	session.Values["csrf"] = token
	err = session.Save(r, w)
	if err != nil {
		core.Warning("Failed to save to session: " + err.Error())
		return err
	}

	http.SetCookie(w, CreateCookie("client-csrf", token, 0, false))
	return nil
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
	return session.Values["csrf"].(string) == token, nil
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
