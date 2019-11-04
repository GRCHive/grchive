package webcore

import (
	"github.com/google/uuid"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"net/http"
	"time"
)

// Returns the raw unhashed API key, the database entry, and an error if any.
func GenerateTemporaryAPIKeyForUser(userId int64) (core.RawApiKey, *core.ApiKey) {
	rawKey := core.RawApiKey(uuid.New().String())
	salt := uuid.New().String()
	hashedKey := rawKey.HashWithSalt(salt)

	key := core.ApiKey{
		HashedKey:      hashedKey,
		Salt:           salt,
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

	if key == nil || key.IsExpired() {
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
		http.SetCookie(w, CreateCookie("client-api-key", string(rawKey), key.SecondsToExpiration(), false))
	}

	return nil
}
