package backblaze

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const AuthApiEndpoint string = "https://api.backblazeb2.com/b2api/v2/b2_authorize_account"

type B2Key struct {
	Id  string
	Key string
}

type B2AuthToken struct {
	Token       string `json:"authorizationToken"`
	Expiration  time.Time
	ApiUrl      string `json:"apiUrl"`
	DownloadUrl string `json:"downloadUrl"`
}

// Maps an application key to an auth token.
// Key: B2Key
// Value: *B2AuthToken
var authTokens = sync.Map{}

func B2Auth(key B2Key) (*B2AuthToken, error) {
	// If token doesn't exist, use the b2_authorize_account
	// API to create a token. If the token is expired,
	// do the same thing.
	// TODO: Do we want to not send multiple requests for the some key
	// 		 if they come within a short amount of time?
	rawToken, ok := authTokens.Load(key)
	var token *B2AuthToken = nil
	if ok {
		token = rawToken.(*B2AuthToken)
	}

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

		if err = handleBackblazeError(resp, rootObj); err != nil {
			return nil, err
		}

		newToken := B2AuthToken{
			// Expiration should be OK for 24 hours but only
			// refresh every 12 hours to be safe.
			Expiration: time.Now().UTC().Add(12 * time.Hour),
		}

		// At this point we should be pretty confident that
		// we are actually authorized so we should be able to unmarshal
		// directly in the B2AuthToken struct now.
		err = json.Unmarshal(respBodyData, &newToken)
		if err != nil {
			return nil, err
		}

		authTokens.Store(key, &newToken)
		token = &newToken
	}

	return token, nil
}
