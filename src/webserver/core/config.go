package core

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

type GlobalConfig struct {
	CompanyName string
	Domain      string
}

var globalConfig *GlobalConfig

type EnvConfig struct {
	DatabaseDriver     string
	DatabaseConnString string
}

var envConfig *EnvConfig

func LoadGlobalProps() *GlobalConfig {
	if globalConfig == nil {
		globalConfig = &GlobalConfig{
			CompanyName: "Audit Stuff",
			Domain:      "auditstuff.com",
		}
	}
	return globalConfig
}

func LoadEnvConfig() *EnvConfig {
	if envConfig == nil {
		dat, err := ioutil.ReadFile("src/webserver/config/config.toml")
		if err != nil {
			Error(err.Error())
		}
		tomlConfig, terr := toml.Load(string(dat))
		if terr != nil {
			Error(err.Error())
		}

		envConfig = new(EnvConfig)
		envConfig.DatabaseDriver = tomlConfig.Get("database.type").(string)
		envConfig.DatabaseConnString = tomlConfig.Get("database.connection").(string)
	}

	return envConfig
}
