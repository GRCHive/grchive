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
	AccessToken  string  `json:"access_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    float64 `json:"expires_in"`
	Scope        string  `json:"scope"`
	RefreshToken string  `json:"refresh_token"`
	IdToken      string  `json:"id_token"`
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

// Creates a core.UserSession object and stores it into the session database.
func CreateUserSesssionFromTokens(tokens *OktaTokens, accessJwt *RawJWT, idJwt *RawJWT, r *http.Request) (*core.UserSession, error) {
	if len(idJwt.Payload.Email) == 0 || len(accessJwt.Payload.Sub) == 0 {
		return nil, errors.New("Failed to find email in ID/Access Token.")
	}

	// Is this even necessary?
	if idJwt.Payload.Email != accessJwt.Payload.Sub {
		return nil, errors.New("Id token Email vs Access Token sub mismatch.")
	}

	userSession := &core.UserSession{
		SessionId:      uuid.New().String(),
		UserEmail:      idJwt.Payload.Email,
		LastActiveTime: time.Now(),
		ExpirationTime: time.Unix(accessJwt.Payload.Exp, 0),
		UserAgent:      r.UserAgent(),
		IP:             r.RemoteAddr,
		AccessToken:    tokens.AccessToken,
		IdToken:        tokens.IdToken,
		RefreshToken:   tokens.RefreshToken,
	}

	err := database.StoreUserSession(userSession)
	if err != nil {
		return nil, err
	}

	return userSession, nil
}

func OktaObtainTokens(code string, r *http.Request) (*core.UserSession, error) {
	envConfig := core.LoadEnvConfig()

	var postVals url.Values = url.Values{
		"code":          []string{code},
		"grant_type":    []string{envConfig.Login.GrantType},
		"redirect_uri":  []string{core.FullSamlCallbackUrl},
		"client_id":     []string{envConfig.Login.ClientId},
		"client_secret": []string{envConfig.Login.ClientSecret},
		"scope":         []string{url.QueryEscape(envConfig.Login.Scope)},
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

	return CreateUserSesssionFromTokens(data, accessJwt, idJwt, r)
}

// If successful, returns a new http.Request that contains
// a context.Context with the user session. Otherwise, returns a nil
// along with an error.
func VerifyUserSessionAuthenticated(r *http.Request) (*http.Request, error) {
	sessionId, err := GetUserSessionOnClient(r)
	if err != nil {
		return nil, err
	}

	session, err := database.FindUserSession(sessionId)
	if err != nil {
		return nil, err
	}

	ctx := AddSessionToContext(session, r.Context())
	return r.Clone(ctx), nil
}
