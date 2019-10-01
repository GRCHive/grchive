package core

import (
	"fmt"
	"net/url"
)

// Landing Page
var GetStartedUrl string = "/getting-started"
var ContactUsUrl string = "/contact-us"
var HomePageUrl string = "/"
var LoginUrl string = "/login"
var LearnMoreUrl string = "/learn"
var SamlCallbackUrl string = LoadEnvConfig().Login.RedirectUrl
var FullSamlCallbackUrl string = LoadEnvConfig().SelfUri + SamlCallbackUrl

// Dashboard
var DashboardUrl string = "/dashboard"

func CreateOktaLoginUrl(idp string, state string, nonce string) string {
	envConfig := LoadEnvConfig()

	return fmt.Sprintf("%s%s?idp=%s&client_id=%s&response_type=%s&response_mode=%s&scope=%s&redirect_uri=%s&state=%s&=nonce=%s",
		envConfig.Login.BaseUrl,
		envConfig.Login.AuthEndpoint,
		idp,
		envConfig.Login.ClientId,
		envConfig.Login.ResponseType,
		envConfig.Login.ResponseMode,
		url.QueryEscape(envConfig.Login.Scope),
		url.QueryEscape(FullSamlCallbackUrl),
		state,
		url.QueryEscape(nonce))
}

var OktaTokenUrl string = fmt.Sprintf("%s%s",
	LoadEnvConfig().Login.BaseUrl,
	LoadEnvConfig().Login.TokenEndpoint)
