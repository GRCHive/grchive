package backblaze

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const AuthApiEndpoint string = "https://api.backblazeb2.com/b2api/v2/b2_authorize_account"

type B2Key struct {
	Id  string
	Key string
}

type B2AuthToken struct {
	Token       string
	Expiration  time.Time
	ApiUrl      string
	DownloadUrl string
}

// Maps an application key to an auth token.
var authTokens = map[B2Key]*B2AuthToken{}

func B2Auth(key B2Key) (*B2AuthToken, error) {
	// If token doesn't exist, use the b2_authorize_account
	// API to create a token. If the token is expired,
	// do the same thing.
	token, ok := authTokens[key]
	if !ok || time.Now().UTC().After(token.Expiration.UTC()) {
		body := &bytes.Buffer{}
		req, err := http.NewRequest("GET", AuthApiEndpoint, body)
		if err != nil {
			return nil, err
		}
		SetAuthorizationHeaderWithAppKey(req, key)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		respBodyData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		rootObj := map[string]*json.RawMessage{}
		err = json.Unmarshal(respBodyData, &rootObj)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(fmt.Sprintf(
				"Failed to auth with Backblaze: %d\n\tMessage: %s - %s",
				resp.StatusCode,
				string(*rootObj["message"]),
				string(*rootObj["code"]),
			))
		}

		newToken := B2AuthToken{
			Token:       string(*rootObj["authorizationToken"]),
			ApiUrl:      string(*rootObj["apiUrl"]),
			DownloadUrl: string(*rootObj["downloadUrl"]),
			// Expiration should be OK for 24 hours but only
			// refresh every 12 hours to be safe.
			Expiration: time.Now().UTC().Add(12 * time.Hour),
		}
		authTokens[key] = &newToken
		token = &newToken
	}

	return token, nil
}
