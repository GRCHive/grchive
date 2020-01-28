package webcore

import (
	"github.com/google/uuid"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"net/http"
	"time"
)

// Returns the raw unhashed API key, the database entry, and an error if any.
func GenerateTemporaryAPIKeyForUser(userId int64) (core.RawApiKey, *core.ApiKey) {
	rawKey := core.RawApiKey(uuid.New().String())
	hashedKey := rawKey.Hash()

	key := core.ApiKey{
		HashedKey:      hashedKey,
		ExpirationDate: time.Now().Add(time.Hour).UTC(),
		UserId:         userId,
	}

	return rawKey, &key
}

func RefreshGrantAPIKey(userId int64, w http.ResponseWriter, r *http.Request) error {
	key, err := database.FindApiKeyForUser(userId)
	if err != nil {
		return err
	}

	// Check if the user's API key is the same as what we have on the server.
	// If it isn't the same, re-grant the API key.
	// This situation happens if they share the browser and login as a different user.
	forceNeedNewKey := false
	cookie, err := r.Cookie("client-api-key")
	if err == nil && key != nil {
		clientKey := core.RawApiKey(cookie.Value)
		forceNeedNewKey = clientKey.Hash() != key.HashedKey
	}

	if key == nil || forceNeedNewKey || key.IsExpired(core.DefaultClock) {
		isNew := (key == nil)
		rawKey, key := GenerateTemporaryAPIKeyForUser(userId)

		if isNew {
			err = database.StoreApiKey(key)
		} else {
			err = database.UpdateApiKey(key)
		}

		if err != nil {
			return err
		}
		http.SetCookie(w, CreateCookie("client-api-key", string(rawKey), key.SecondsToExpiration(core.DefaultClock), false))
	}

	return nil
}

func GetAPIKeyFromRequest(r *http.Request) (*core.ApiKey, error) {
	rawApiKey := GetRawClientAPIKeyFromRequest(r)
	hashedRawKey := rawApiKey.Hash()
	return database.FindApiKey(hashedRawKey)
}
