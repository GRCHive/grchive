package core

import (
	"fmt"
	"net/url"
)

var GetStartedUrl string = "/getting-started"
var ContactUsUrl string = "/contact-us"
var HomePageUrl string = "/"
var LoginUrl string = "/login"
var LearnMoreUrl string = "/learn"

func CreateOktaLoginUrl(idp string, state string, nonce string) string {
	envConfig := LoadEnvConfig()

	return fmt.Sprintf("%s?idp=%s&client_id=%s&response_type=%s&response_mode=%s&scope=%s&redirect_uri=%s&state=%s&=nonce=%s",
		envConfig.Login.BaseUrl,
		idp,
		envConfig.Login.ClientId,
		envConfig.Login.ResponseType,
		envConfig.Login.ResponseMode,
		envConfig.Login.Scope,
		url.QueryEscape(envConfig.Login.RedirectUrl),
		state,
		url.QueryEscape(nonce))
}
