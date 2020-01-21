package core

import (
	"encoding/hex"
	"github.com/pelletier/go-toml"
	"gitlab.com/b3h47pte/audit-stuff/backblaze_api"
	"gitlab.com/b3h47pte/audit-stuff/mail_api"
	"io/ioutil"
)

var TomlConfig *toml.Tree

type CompanyConfig struct {
	CompanyName string
	Domain      string
}

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

type VaultConfig struct {
	Url   string
	Token string
}

type BackblazeConfig struct {
	Key                backblaze.B2Key
	ControlDocBucketId string
}

type MailConfig struct {
	Provider      mail.MailAPIProvider
	Key           string
	VeriEmailFrom mail.Email
}

type HashIdConfigData struct {
	Salt      string
	MinLength int
}

type RabbitMQConfig struct {
	Host string
	Port int32
}

type EnvConfigData struct {
	SelfUri            string
	DatabaseConnString string
	Login              *LoginConfig
	Okta               *OktaConfig
	SessionKeys        [][]byte
	UseSecureCookies   bool
	Company            *CompanyConfig
	Vault              *VaultConfig
	Backblaze          *BackblazeConfig
	Mail               *MailConfig
	HashId             *HashIdConfigData
	RabbitMQ           *RabbitMQConfig
}

var EnvConfig *EnvConfigData

// I think there might be a Bazel bug here where
// bazel run puts us into the runfile/workspace_name folder
// instead of the runfile folder?
var DefaultConfigLocation = "src/webserver/config/config.toml"

func LoadTomlConfig(configLoc string) *toml.Tree {
	dat, err := ioutil.ReadFile(configLoc)
	if err != nil {
		Error(err.Error())
	}
	tomlConfig, err := toml.Load(string(dat))
	if err != nil {
		Error(err.Error())
	}
	return tomlConfig
}

func LoadEnvConfig(tomlConfig *toml.Tree) *EnvConfigData {
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

	envConfig.Company = new(CompanyConfig)
	envConfig.Company.CompanyName = tomlConfig.Get("company.company_name").(string)
	envConfig.Company.Domain = tomlConfig.Get("company.domain").(string)

	envConfig.Vault = new(VaultConfig)
	envConfig.Vault.Url = tomlConfig.Get("vault.url").(string)
	envConfig.Vault.Token = tomlConfig.Get("vault.token").(string)

	envConfig.Backblaze = new(BackblazeConfig)
	envConfig.Backblaze.ControlDocBucketId = tomlConfig.Get("backblaze.control_doc_bucket").(string)
	envConfig.Backblaze.Key.Id = tomlConfig.Get("backblaze.keyId").(string)
	envConfig.Backblaze.Key.Key = tomlConfig.Get("backblaze.key").(string)

	envConfig.Mail = new(MailConfig)
	envConfig.Mail.Provider = mail.MailAPIProvider(tomlConfig.Get("mail.provider").(string))
	envConfig.Mail.Key = tomlConfig.Get("mail.key").(string)
	envConfig.Mail.VeriEmailFrom.Name = tomlConfig.Get("mail.verification.from.name").(string)
	envConfig.Mail.VeriEmailFrom.Email = tomlConfig.Get("mail.verification.from.email").(string)

	envConfig.HashId = new(HashIdConfigData)
	envConfig.HashId.Salt = tomlConfig.Get("hashids.salt").(string)
	envConfig.HashId.MinLength = int(tomlConfig.Get("hashids.min_length").(int64))

	envConfig.RabbitMQ = new(RabbitMQConfig)
	envConfig.RabbitMQ.Host = tomlConfig.Get("rabbitmq.host").(string)
	envConfig.RabbitMQ.Port = int32(tomlConfig.Get("rabbitmq.port").(int64))

	return envConfig
}

func InitializeConfig(configLoc string) {
	TomlConfig = LoadTomlConfig(configLoc)
	EnvConfig = LoadEnvConfig(TomlConfig)
}