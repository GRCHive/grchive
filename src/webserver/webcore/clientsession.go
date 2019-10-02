package webcore

import (
	"github.com/gorilla/Sessions"
	"github.com/gorilla/securecookie"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"net/http"
	"time"
)

var ClientShortSessionStore = sessions.NewCookieStore(core.LoadEnvConfig().SessionKeys...)
var ClientLongSessionStore = sessions.NewCookieStore(core.LoadEnvConfig().SessionKeys...)

// The first object in the slice is used to encrypt all cookies. The other objects are there
// to handle the case of key rotation.
var Cookies = make([]*securecookie.SecureCookie, 0)

func InitializeSessions() {
	ClientShortSessionStore.Options.HttpOnly = true
	ClientShortSessionStore.Options.Secure = core.LoadEnvConfig().UseSecureCookies
	ClientShortSessionStore.Options.MaxAge = 0

	ClientLongSessionStore.Options.HttpOnly = true
	ClientLongSessionStore.Options.Secure = core.LoadEnvConfig().UseSecureCookies
	ClientLongSessionStore.Options.MaxAge = core.SecondsInDay * 30

	// Create a new SecureCookie object for every session key pair.
	encryptKeys := core.LoadEnvConfig().SessionKeys
	for i := 0; i < len(encryptKeys)/2; i++ {
		hashKey := encryptKeys[i*2]
		blockKey := encryptKeys[i*2+1]
		Cookies = append(Cookies, securecookie.New(hashKey, blockKey))
	}
}

func StoreUserSessionOnClient(session *core.UserSession, w http.ResponseWriter) error {
	// Bypass gorilla sessions and use gorilla securecookie directly to create the cookie ourselves
	// so that we can craft the expiration of the cookie to match the expiration of the session on
	// the server.
	value := map[string]string{
		"sessionId": session.SessionId,
	}

	cookieName := "userSession"

	var encoded string
	var err error
	if encoded, err = Cookies[0].Encode(cookieName, value); err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    encoded,
		Expires:  session.ExpirationTime,
		MaxAge:   int(session.ExpirationTime.Sub(time.Now()).Seconds()),
		Secure:   core.LoadEnvConfig().UseSecureCookies,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return nil
}