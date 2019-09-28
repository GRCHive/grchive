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

type EnvConfig struct {
	DatabaseConnString string
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
	}

	return envConfig
}
