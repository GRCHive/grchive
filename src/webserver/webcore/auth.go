package webcore

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"time"
)

type OktaKey struct {
	Status string `json:"status"`
	Kid    string `json:"kid"`
	N      string `json:"n"`
	E      string `json:"e"`
}

type OktaTokens struct {
	AccessToken        string  `json:"access_token"`
	TokenType          string  `json:"token_type"`
	ExpiresIn          float64 `json:"expires_in"`
	Scope              string  `json:"scope"`
	RefreshToken       string  `json:"refresh_token"`
	IdToken            string  `json:"id_token"`
	DecodedAccessToken *RawJWT
	DecodedIDToken     *RawJWT
}

type OktaJWTRetriever struct{}

var oktaJwtRetriever *OktaJWTRetriever = new(OktaJWTRetriever)
var oktaJwtManager *JWTManager = &JWTManager{impl: oktaJwtRetriever}

func (this OktaJWTRetriever) RetrieveKeys() (map[string][]*rsa.PublicKey, error) {
	resp, err := http.Get(core.OktaKeyUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rootObj := map[string]*json.RawMessage{}
	err = json.Unmarshal(body, &rootObj)
	if err != nil {
		return nil, err
	}

	if val, ok := rootObj["keys"]; ok {
		allKeys := make([]OktaKey, 0)
		err = json.Unmarshal(*val, &allKeys)
		if err != nil {
			return nil, err
		}

		retMap := make(map[string][]*rsa.PublicKey)
		for i := 0; i < len(allKeys); i++ {
			if allKeys[i].Status == "EXPIRED" {
				continue
			}

			// See https://tools.ietf.org/html/rfc7518#page-30
			// N and E are Base64urlUint encoded.
			byteN, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(allKeys[i].N)
			if err != nil {
				return nil, err
			}

			byteE, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(allKeys[i].E)
			if err != nil {
				return nil, err
			}

			var N, E big.Int
			N.SetBytes(byteN)
			E.SetBytes(byteE)

			newKey := &rsa.PublicKey{
				N: &N,
				E: int(E.Int64()),
			}

			if _, ok := retMap[allKeys[i].Kid]; ok {
				retMap[allKeys[i].Kid] = append(retMap[allKeys[i].Kid], newKey)
			} else {
				retMap[allKeys[i].Kid] = []*rsa.PublicKey{newKey}
			}
		}

		return retMap, nil
	}
	return nil, errors.New("Failed to find 'keys' key in okta retrieve keys.")
}

func UpdateUserSessionFromTokens(session *core.UserSession, tokens *OktaTokens, r *http.Request) error {
	accessJwt := tokens.DecodedAccessToken
	idJwt := tokens.DecodedIDToken

	if len(idJwt.Payload.Email) == 0 || len(accessJwt.Payload.Sub) == 0 {
		return errors.New("Failed to find email in ID/Access Token.")
	}

	// Is this even necessary?
	if idJwt.Payload.Email != accessJwt.Payload.Sub {
		return errors.New("Id token Email vs Access Token sub mismatch.")
	}

	*session = core.UserSession{
		SessionId:      session.SessionId,
		UserEmail:      idJwt.Payload.Email,
		LastActiveTime: time.Now().UTC(),
		ExpirationTime: time.Unix(accessJwt.Payload.Exp, 0).UTC(),
		UserAgent:      r.UserAgent(),
		IP:             r.RemoteAddr,
		AccessToken:    tokens.AccessToken,
		IdToken:        tokens.IdToken,
		RefreshToken:   tokens.RefreshToken,
	}

	// Create new session ID if it's empty. This allows for updated sessions (aka
	// session with an existing session ID to not receive a new ID).
	if len(session.SessionId) == 0 {
		session.SessionId = uuid.New().String()
	}

	return nil
}

// Creates a core.UserSession object and stores it into the session database.
func CreateUserSessionFromTokens(tokens *OktaTokens, r *http.Request) (*core.UserSession, error) {
	userSession := new(core.UserSession)
	err := UpdateUserSessionFromTokens(userSession, tokens, r)
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func OktaObtainTokens(code string, isRefresh bool) (*OktaTokens, error) {
	envConfig := core.LoadEnvConfig()

	var postVals url.Values = url.Values{
		"redirect_uri":  []string{core.FullSamlCallbackUrl},
		"client_id":     []string{envConfig.Login.ClientId},
		"client_secret": []string{envConfig.Login.ClientSecret},
		"scope":         []string{envConfig.Login.Scope},
	}

	if isRefresh {
		postVals["refresh_token"] = []string{code}
		postVals["grant_type"] = []string{"refresh_token"}
	} else {
		postVals["code"] = []string{code}
		postVals["grant_type"] = []string{envConfig.Login.GrantType}
	}

	resp, err := http.PostForm(core.OktaTokenUrl, postVals)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data *OktaTokens = new(OktaTokens)
	err = json.Unmarshal(body, data)
	if err != nil || data == nil {
		return nil, err
	}

	accessJwt, err := oktaJwtManager.VerifyJWT(data.AccessToken, true)
	if err != nil {
		return nil, err
	}

	idJwt, err := oktaJwtManager.VerifyJWT(data.IdToken, false)
	if err != nil {
		return nil, err
	}

	data.DecodedAccessToken = accessJwt
	data.DecodedIDToken = idJwt
	return data, nil
}

// If successful, returns a new http.Request that contains
// a context.Context with the user session. Otherwise, returns the passed
// in request along with an error. If a session was found and it can't
// be validated, it will be deleted. If a session was found and it is
// validated, its last active time will be updated and its access token
// will be refreshed if necessary. Also update the user session cookie's
// expiration time.
func FindValidUserSession(w http.ResponseWriter, r *http.Request) (*core.UserSession, *http.Request, error) {
	sessionId, err := GetUserSessionOnClient(r)
	if err != nil {
		return nil, r, err
	}

	session, err := database.FindUserSession(sessionId)
	if err != nil {
		return nil, r, err
	}

	// Ensure that the access token is valid and not expired. If it is,
	// use the corresponding refresh token to retrieve a new token.
	// If this fails, force the user to re-login. Note that we have three
	// sources where we check for expiration: 1) access token 2) id token and
	// 3) our session database. All three must not be expired.
	_, accessErr := oktaJwtManager.VerifyJWT(session.AccessToken, true)
	_, idErr := oktaJwtManager.VerifyJWT(session.IdToken, false)
	if core.IsPastTime(session.ExpirationTime) || idErr == ExpiredJWTToken || accessErr == ExpiredJWTToken || true {
		newTokens, err := OktaObtainTokens(session.RefreshToken, true)
		if err != nil {
			return nil, r, err
		}

		err = UpdateUserSessionFromTokens(session, newTokens, r)
		if err != nil {
			return nil, r, err
		}
	} else if accessErr != nil {
		return nil, r, accessErr
	} else if idErr != nil {
		return nil, r, idErr
	}

	// Update last active time and save in DB.
	// This needs to be here because the last active time needs to be
	// valid even if we don't refresh the token.
	session.LastActiveTime = time.Now().UTC()
	// At this point all changes to the session should be finished and
	// any errors should hopefully not point to an invalid session.

	ctx := AddSessionToContext(session, r.Context())
	newR := r.Clone(ctx)

	err = database.UpdateUserSession(session)
	if err != nil {
		// This is probably our fault so don't delete the session but keep it around
		// until the next user request and hopefully it'll resolve itself.
		return session, newR, err
	}

	// Re-store the cookie on the user side.
	err = StoreUserSessionOnClient(session, w)
	if err != nil {
		// Something went wrong with storing the session cookie but
		// we still have a valid cookie/session so just go ahead and
		// use the session.
		return session, newR, err
	}

	return session, newR, nil
}
