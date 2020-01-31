package vault

import (
	"encoding/json"
	"errors"
)

func userPassAuth(username string, password string) error {
	data, err := sendVaultRequest(
		"POST",
		"/v1/auth/userpass/login/"+username,
		struct {
			Password string `json:"password"`
		}{
			Password: password,
		})
	if err != nil {
		return err
	}

	rawAuth, ok := data["auth"]
	if !ok {
		return errors.New("Failed to find auth object.")
	}

	auth := map[string]interface{}{}
	err = json.Unmarshal(*rawAuth, &auth)
	if err != nil {
		return err
	}

	rawToken, ok := auth["client_token"]
	if !ok {
		return errors.New("Failed to find token object.")
	}

	GlobalConfig.token = rawToken.(string)
	return nil
}
