package webcore

import (
	"encoding/base64"
	"encoding/json"
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
	if err != nil {
		forceNeedNewKey = true
	}

	if !forceNeedNewKey {
		var clientApiKey core.RawApiKey = ""

		// Check if cookie is well formed.
		// We expect the cookie to be base64 encoded JSON.
		// If it isn't, then re-grant the cookie (older API keys don't have this form).
		rawJsonCookie, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err == nil {
			jsonData := map[string]string{}
			err := json.Unmarshal([]byte(rawJsonCookie), &jsonData)
			if err == nil {
				clientApiKey = core.RawApiKey(jsonData["Key"])
			} else {
				forceNeedNewKey = true
			}
		} else {
			forceNeedNewKey = true
		}

		if key != nil && !forceNeedNewKey {
			forceNeedNewKey = clientApiKey.Hash() != key.HashedKey
		}
	}

	if err != nil || key == nil || forceNeedNewKey || key.NeedsRefresh(core.DefaultClock) {
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

		data, err := json.Marshal(struct {
			Key        string
			Expiration time.Time
		}{
			Key:        string(rawKey),
			Expiration: key.ExpirationDate,
		})

		if err != nil {
			return err
		}
		http.SetCookie(w, CreateCookie("client-api-key", base64.StdEncoding.EncodeToString(data), key.SecondsToExpiration(core.DefaultClock), false))
	}

	return nil
}

func GetAPIKeyFromRequest(w http.ResponseWriter, r *http.Request) (*core.ApiKey, error) {
	rawApiKey := GetRawClientAPIKeyFromRequest(r)
	hashedRawKey := rawApiKey.Hash()
	key, err := database.FindApiKey(hashedRawKey)
	if err != nil {
		return nil, err
	}

	if key.IsExpired(core.DefaultClock) {
		// For convenience we need to regenerate the API key.
		// TODO: Protect ourselves against replay attacks?
		err = RefreshGrantAPIKey(key.UserId, w, r)
		if err != nil {
			return nil, err
		}
		return database.FindApiKeyForUser(key.UserId)
	}
	return key, nil
}
