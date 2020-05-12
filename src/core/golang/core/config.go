package core

import (
	"encoding/hex"
	"github.com/pelletier/go-toml"
	"gitlab.com/grchive/grchive/mail_api"
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
	Url      string
	Username string
	Password string
}

type GCloudConfig struct {
	AuthFilename string
	DocBucket    string
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
	Username string
	Password string
	Host     string
	Port     int32
	UseTLS   bool
}

type GrpcConfig struct {
	Host    string
	Port    int32
	TLS     bool
	TLSCert string
	TLSKey  string
}

type GrpcEndpoints struct {
	QueryRunner GrpcConfig
}

type TLSConfig struct {
	TLSRootCaCert string
}

type ContainerRegistryAuth struct {
	Username string
	Password string
}

type KotlinConfigData struct {
	GroupId      string
	ArtifactId   string
	MajorVersion int
	MinorVersion int
}

type GiteaConfigData struct {
	Token     string
	Host      string
	Port      int32
	Protocol  string
	GlobalOrg string
}

type DroneConfigData struct {
	Token                   string
	Host                    string
	Port                    int32
	Protocol                string
	RunnerType              string
	RunnerImage             string
	RunnerImagePull         string
	RunnerDbConnectOverride string
}

type FeatureFlags struct {
	Automation bool
}

type ArtifactoryConfigData struct {
	Host string
	Port int32
}

type NotificationConfig struct {
	EnableEmail bool
}

type ScriptRunnerConfig struct {
	RunnerImage string
}

type EnvConfigData struct {
	SelfUri            string
	SelfDomain         string
	DatabaseConnString string
	DatabaseUsername   string
	DatabasePassword   string
	Login              *LoginConfig
	Okta               *OktaConfig
	SessionKeys        [][]byte
	HmacKey            []byte
	UseSecureCookies   bool
	LogEncryptionPath  string
	Company            *CompanyConfig
	Vault              *VaultConfig
	Gcloud             *GCloudConfig
	Mail               *MailConfig
	HashId             *HashIdConfigData
	RabbitMQ           *RabbitMQConfig
	Grpc               *GrpcEndpoints
	Tls                *TLSConfig
	GitlabRegistryAuth *ContainerRegistryAuth
	Kotlin             *KotlinConfigData
	Gitea              *GiteaConfigData
	Drone              *DroneConfigData
	Features           *FeatureFlags
	Artifactory        *ArtifactoryConfigData
	Notifications      *NotificationConfig
	ScriptRunner       ScriptRunnerConfig
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
	envConfig.SelfDomain = tomlConfig.Get("self_domain").(string)
	envConfig.DatabaseConnString = tomlConfig.Get("database.connection").(string)
	envConfig.DatabaseUsername = tomlConfig.Get("database.username").(string)
	envConfig.DatabasePassword = tomlConfig.Get("database.password").(string)

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

	envConfig.HmacKey, err = hex.DecodeString(tomlConfig.Get("security.hmac_key").(string))
	if err != nil {
		Error(err.Error())
	}

	envConfig.UseSecureCookies = tomlConfig.Get("security.use_secure_cookies").(bool)
	envConfig.LogEncryptionPath = tomlConfig.Get("security.log_encryption_path").(string)

	envConfig.Company = new(CompanyConfig)
	envConfig.Company.CompanyName = tomlConfig.Get("company.company_name").(string)
	envConfig.Company.Domain = tomlConfig.Get("company.domain").(string)

	envConfig.Vault = new(VaultConfig)
	envConfig.Vault.Url = tomlConfig.Get("vault.url").(string)
	envConfig.Vault.Username = tomlConfig.Get("vault.userpass.username").(string)
	envConfig.Vault.Password = tomlConfig.Get("vault.userpass.password").(string)

	envConfig.Gcloud = new(GCloudConfig)
	envConfig.Gcloud.AuthFilename = tomlConfig.Get("gcloud.credentials_file").(string)
	envConfig.Gcloud.DocBucket = tomlConfig.Get("gcloud.storage.doc_bucket").(string)

	envConfig.Mail = new(MailConfig)
	envConfig.Mail.Provider = mail.MailAPIProvider(tomlConfig.Get("mail.provider").(string))
	envConfig.Mail.Key = tomlConfig.Get("mail.key").(string)
	envConfig.Mail.VeriEmailFrom.Name = tomlConfig.Get("mail.verification.from.name").(string)
	envConfig.Mail.VeriEmailFrom.Email = tomlConfig.Get("mail.verification.from.email").(string)

	envConfig.HashId = new(HashIdConfigData)
	envConfig.HashId.Salt = tomlConfig.Get("hashids.salt").(string)
	envConfig.HashId.MinLength = int(tomlConfig.Get("hashids.min_length").(int64))

	envConfig.RabbitMQ = new(RabbitMQConfig)
	envConfig.RabbitMQ.Username = tomlConfig.Get("rabbitmq.username").(string)
	envConfig.RabbitMQ.Password = tomlConfig.Get("rabbitmq.password").(string)
	envConfig.RabbitMQ.Host = tomlConfig.Get("rabbitmq.host").(string)
	envConfig.RabbitMQ.Port = int32(tomlConfig.Get("rabbitmq.port").(int64))
	envConfig.RabbitMQ.UseTLS = tomlConfig.Get("rabbitmq.use_tls").(bool)

	envConfig.Grpc = new(GrpcEndpoints)
	envConfig.Grpc.QueryRunner.Host = tomlConfig.Get("grpc.query_runner.host").(string)
	envConfig.Grpc.QueryRunner.Port = int32(tomlConfig.Get("grpc.query_runner.port").(int64))
	envConfig.Grpc.QueryRunner.TLS = tomlConfig.Get("grpc.query_runner.tls.enable").(bool)
	envConfig.Grpc.QueryRunner.TLSCert = tomlConfig.Get("grpc.query_runner.tls.cert").(string)
	envConfig.Grpc.QueryRunner.TLSKey = tomlConfig.Get("grpc.query_runner.tls.key").(string)

	envConfig.Tls = new(TLSConfig)
	envConfig.Tls.TLSRootCaCert = tomlConfig.Get("tls.root_ca").(string)

	envConfig.GitlabRegistryAuth = new(ContainerRegistryAuth)
	envConfig.GitlabRegistryAuth.Username = tomlConfig.Get("registry.gitlab.username").(string)
	envConfig.GitlabRegistryAuth.Password = tomlConfig.Get("registry.gitlab.password").(string)

	envConfig.Kotlin = new(KotlinConfigData)
	envConfig.Kotlin.GroupId = tomlConfig.Get("kotlin.group_id").(string)
	envConfig.Kotlin.ArtifactId = tomlConfig.Get("kotlin.artifact_id").(string)
	envConfig.Kotlin.MajorVersion = int(tomlConfig.Get("kotlin.major_version").(int64))
	envConfig.Kotlin.MinorVersion = int(tomlConfig.Get("kotlin.minor_version").(int64))

	envConfig.Gitea = new(GiteaConfigData)
	envConfig.Gitea.Host = tomlConfig.Get("gitea.host").(string)
	envConfig.Gitea.Port = int32(tomlConfig.Get("gitea.port").(int64))
	envConfig.Gitea.Protocol = tomlConfig.Get("gitea.protocol").(string)
	envConfig.Gitea.Token = tomlConfig.Get("gitea.token").(string)
	envConfig.Gitea.GlobalOrg = tomlConfig.Get("gitea.global_org").(string)

	envConfig.Drone = new(DroneConfigData)
	envConfig.Drone.Host = tomlConfig.Get("drone.host").(string)
	envConfig.Drone.Port = int32(tomlConfig.Get("drone.port").(int64))
	envConfig.Drone.Protocol = tomlConfig.Get("drone.protocol").(string)
	envConfig.Drone.Token = tomlConfig.Get("drone.token").(string)
	envConfig.Drone.RunnerType = tomlConfig.Get("drone.runner_type").(string)
	envConfig.Drone.RunnerImage = tomlConfig.Get("drone.runner_image").(string)
	envConfig.Drone.RunnerImagePull = tomlConfig.Get("drone.runner_image_pull").(string)
	envConfig.Drone.RunnerDbConnectOverride = tomlConfig.Get("drone.runner_db_connect_override").(string)

	envConfig.ScriptRunner.RunnerImage = tomlConfig.Get("script_runner.runner_image").(string)

	envConfig.Features = new(FeatureFlags)
	envConfig.Features.Automation = tomlConfig.Get("features.automation").(bool)

	envConfig.Artifactory = new(ArtifactoryConfigData)
	envConfig.Artifactory.Host = tomlConfig.Get("artifactory.host").(string)
	envConfig.Artifactory.Port = int32(tomlConfig.Get("artifactory.port").(int64))

	envConfig.Notifications = new(NotificationConfig)
	envConfig.Notifications.EnableEmail = tomlConfig.Get("notifications.enable_email").(bool)

	return envConfig
}

func InitializeConfig(configLoc string) {
	TomlConfig = LoadTomlConfig(configLoc)
	EnvConfig = LoadEnvConfig(TomlConfig)
}
