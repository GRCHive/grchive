package okta

import (
	"encoding/json"
	"errors"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
)

const RegisterUserEndpoint string = "/api/v1/users"

type ProfileData struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type PasswordData struct {
	Value string `json:"value"`
}

type CredentialData struct {
	Password PasswordData `json:"password"`
}

type RegisterUserData struct {
	Profile     ProfileData    `json:"profile" url:"-"`
	Credentials CredentialData `json:"credentials" url:"-"`
	Activate    bool           `json:"-" url:"activate"`
}

// Returns the Okta User ID
func RegisterUser(user RegisterUserData) (string, error) {
	queryVals, err := query.Values(user)
	if err != nil {
		return "", err
	}

	resp, err := oktaPost(RegisterUserEndpoint+"?"+queryVals.Encode(), user)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Okta Register User failed: " + string(body))
	}

	rootObj := map[string]*json.RawMessage{}
	err = json.Unmarshal(body, &rootObj)
	if err != nil {
		return "", err
	}

	id := ""
	err = json.Unmarshal(*rootObj["id"], &id)
	if err != nil {
		return "", err
	}

	return id, nil
}
