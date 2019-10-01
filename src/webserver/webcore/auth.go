package webcore

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"io/ioutil"
	"net/http"
	"net/url"
)

type OktaTokens struct {
	AccessToken  string  `json: "access_token"`
	TokenType    string  `json: "token_type"`
	ExpiresIn    float64 `json: "expires_in"`
	Scope        string  `json: "scope"`
	RefreshToken string  `json: "refresh_token"`
	IdToken      string  `json: "id_token"`
}

func CreateUserSesssionFromTokens(tokens *OktaTokens, r *http.Request) (*core.UserSession, error) {
	return session, nil
}

func OktaObtainTokens(code string, r *http.Request) (*core.UserSession, error) {
	envConfig := core.LoadEnvConfig()

	var postVals url.Values = url.Values{
		"code":          []string{code},
		"grant_type":    []string{envConfig.Login.GrantType},
		"redirect_uri":  []string{core.FullSamlCallbackUrl},
		"client_id":     []string{envConfig.Login.ClientId},
		"client_secret": []string{envConfig.Login.ClientSecret},
		"scope":         []string{"openid"},
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

	return CreateUserSesssionFromRequest("", r)
}
