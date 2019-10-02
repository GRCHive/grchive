package core

import (
	"encoding/hex"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

var tomlConfig *toml.Tree

type TemplateConfig struct {
	CompanyName  string
	Domain       string
	RecaptchaKey string
}

var templateConfig *TemplateConfig

type LoginConfig struct {
	BaseUrl       string
	AuthEndpoint  string
	TokenEndpoint string
	KeyEndpoint   string
	ClientId      string
	ClientSecret  string
	ResponseType  string
	ResponseMode  string
	Scope         string
	RedirectUrl   string
	GrantType     string
}

type EnvConfig struct {
	SelfUri            string
	DatabaseConnString string
	Login              *LoginConfig
	SessionKeys        [][]byte
	UseSecureCookies   bool
}

var envConfig *EnvConfig

func loadTomlConfig() {
	if tomlConfig == nil {
		dat, err := ioutil.ReadFile("src/webserver/config/config.toml")
		if err != nil {
			Error(err.Error())
		}
		tomlConfig, err = toml.Load(string(dat))
		if err != nil {
			Error(err.Error())
		}
	}
}

func LoadTemplateConfig() *TemplateConfig {
	loadTomlConfig()
	if templateConfig == nil {
		templateConfig = &TemplateConfig{
			CompanyName: "Audit Stuff",
			Domain:      "auditstuff.com",
		}
	}
	return templateConfig
}

func LoadEnvConfig() *EnvConfig {
	var err error

	loadTomlConfig()
	if envConfig == nil {
		envConfig = new(EnvConfig)
		envConfig.SelfUri = tomlConfig.Get("self_uri").(string)
		envConfig.DatabaseConnString = tomlConfig.Get("database.connection").(string)

		envConfig.Login = new(LoginConfig)
		envConfig.Login.BaseUrl = tomlConfig.Get("login.url").(string)
		envConfig.Login.AuthEndpoint = tomlConfig.Get("login.auth_endpoint").(string)
		envConfig.Login.TokenEndpoint = tomlConfig.Get("login.token_endpoint").(string)
		envConfig.Login.KeyEndpoint = tomlConfig.Get("login.key_endpoint").(string)
		envConfig.Login.ClientId = tomlConfig.Get("login.params.client_id").(string)
		envConfig.Login.ClientSecret = tomlConfig.Get("login.params.client_secret").(string)
		envConfig.Login.ResponseType = tomlConfig.Get("login.params.response_type").(string)
		envConfig.Login.ResponseMode = tomlConfig.Get("login.params.response_mode").(string)
		envConfig.Login.Scope = tomlConfig.Get("login.params.scope").(string)
		envConfig.Login.RedirectUrl = tomlConfig.Get("login.params.redirect_uri").(string)
		envConfig.Login.GrantType = tomlConfig.Get("login.params.grant_type").(string)

		tmpSessionKeys := tomlConfig.Get("security.session_keys").([]interface{})
		envConfig.SessionKeys = make([][]byte, len(tmpSessionKeys))
		for i := 0; i < len(tmpSessionKeys); i++ {
			envConfig.SessionKeys[i], err = hex.DecodeString(tmpSessionKeys[i].(string))
			if err != nil {
				Error(err.Error())
			}
		}
		envConfig.UseSecureCookies = tomlConfig.Get("security.use_secure_cookies").(bool)
	}

	return envConfig
}
