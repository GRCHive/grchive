package core

import (
	"encoding/hex"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

var TomlConfig *toml.Tree

type TemplateConfigData struct {
	CompanyName  string
	Domain       string
	RecaptchaKey string
}

var TemplateConfig *TemplateConfigData

type LoginConfig struct {
	AuthServerEndpoint     string
	AuthEndpoint           string
	TokenEndpoint          string
	KeyEndpoint            string
	LogoutEndpoint         string
	ClientId               string
	ClientSecret           string
	ResponseType           string
	ResponseMode           string
	Scope                  string
	RedirectUrl            string
	GrantType              string
	AuthAudience           string
	TimeDriftLeewaySeconds float64
}

type OktaConfig struct {
	BaseUrl       string
	ApiEndpoint   string
	ApiKey        string
	UsersEndpoint string
}

type EnvConfigData struct {
	SelfUri            string
	DatabaseConnString string
	Login              *LoginConfig
	Okta               *OktaConfig
	SessionKeys        [][]byte
	UseSecureCookies   bool
}

var EnvConfig *EnvConfigData

func loadTomlConfig() *toml.Tree {
	dat, err := ioutil.ReadFile("src/webserver/config/config.toml")
	if err != nil {
		Error(err.Error())
	}
	tomlConfig, err := toml.Load(string(dat))
	if err != nil {
		Error(err.Error())
	}
	return tomlConfig
}

func loadTemplateConfig() *TemplateConfigData {
	templateConfig := &TemplateConfigData{
		CompanyName: "Audit Stuff",
		Domain:      "auditstuff.com",
	}
	return templateConfig
}

func loadEnvConfig(tomlConfig *toml.Tree) *EnvConfigData {
	var err error

	envConfig := new(EnvConfigData)
	envConfig.SelfUri = tomlConfig.Get("self_uri").(string)
	envConfig.DatabaseConnString = tomlConfig.Get("database.connection").(string)

	envConfig.Okta = new(OktaConfig)
	envConfig.Okta.BaseUrl = tomlConfig.Get("okta.url").(string)
	envConfig.Okta.ApiEndpoint = tomlConfig.Get("okta.api_endpoint").(string)
	envConfig.Okta.ApiKey = tomlConfig.Get("okta.api_key").(string)
	envConfig.Okta.UsersEndpoint = tomlConfig.Get("okta.users_endpoint").(string)

	envConfig.Login = new(LoginConfig)
	envConfig.Login.AuthServerEndpoint = tomlConfig.Get("login.authserver_endpoint").(string)
	envConfig.Login.AuthEndpoint = tomlConfig.Get("login.auth_endpoint").(string)
	envConfig.Login.TokenEndpoint = tomlConfig.Get("login.token_endpoint").(string)
	envConfig.Login.KeyEndpoint = tomlConfig.Get("login.key_endpoint").(string)
	envConfig.Login.LogoutEndpoint = tomlConfig.Get("login.logout_endpoint").(string)
	envConfig.Login.ClientId = tomlConfig.Get("login.params.client_id").(string)
	envConfig.Login.ClientSecret = tomlConfig.Get("login.params.client_secret").(string)
	envConfig.Login.ResponseType = tomlConfig.Get("login.params.response_type").(string)
	envConfig.Login.ResponseMode = tomlConfig.Get("login.params.response_mode").(string)
	envConfig.Login.Scope = tomlConfig.Get("login.params.scope").(string)
	envConfig.Login.RedirectUrl = tomlConfig.Get("login.params.redirect_uri").(string)
	envConfig.Login.GrantType = tomlConfig.Get("login.params.grant_type").(string)
	envConfig.Login.AuthAudience = tomlConfig.Get("login.auth_audience").(string)
	envConfig.Login.TimeDriftLeewaySeconds = float64(tomlConfig.Get("login.time_drift_leeway_seconds").(int64))

	tmpSessionKeys := tomlConfig.Get("security.session_keys").([]interface{})
	envConfig.SessionKeys = make([][]byte, len(tmpSessionKeys))
	for i := 0; i < len(tmpSessionKeys); i++ {
		envConfig.SessionKeys[i], err = hex.DecodeString(tmpSessionKeys[i].(string))
		if err != nil {
			Error(err.Error())
		}
	}
	envConfig.UseSecureCookies = tomlConfig.Get("security.use_secure_cookies").(bool)
	return envConfig
}

func InitializeConfig() {
	TomlConfig = loadTomlConfig()
	TemplateConfig = loadTemplateConfig()
	EnvConfig = loadEnvConfig(TomlConfig)
}
