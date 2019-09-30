package core

import (
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
	BaseUrl      string
	ClientId     string
	ResponseType string
	ResponseMode string
	Scope        string
	RedirectUrl  string
}

type EnvConfig struct {
	DatabaseConnString string
	Login              *LoginConfig
	SessionKeys        [][]byte
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
	loadTomlConfig()
	if envConfig == nil {
		envConfig = new(EnvConfig)
		envConfig.DatabaseConnString = tomlConfig.Get("database.connection").(string)

		envConfig.Login = new(LoginConfig)
		envConfig.Login.BaseUrl = tomlConfig.Get("login.url").(string)
		envConfig.Login.ClientId = tomlConfig.Get("login.params.client_id").(string)
		envConfig.Login.ResponseType = tomlConfig.Get("login.params.response_type").(string)
		envConfig.Login.ResponseMode = tomlConfig.Get("login.params.response_mode").(string)
		envConfig.Login.Scope = tomlConfig.Get("login.params.scope").(string)
		envConfig.Login.RedirectUrl = tomlConfig.Get("login.params.redirect_uri").(string)

		tmpSessionKeys := tomlConfig.Get("security.session_keys").([]interface{})
		envConfig.SessionKeys = make([][]byte, len(tmpSessionKeys))
		for i := 0; i < len(tmpSessionKeys); i++ {
			envConfig.SessionKeys[i] = []byte(tmpSessionKeys[i].(string))
		}
	}

	return envConfig
}
