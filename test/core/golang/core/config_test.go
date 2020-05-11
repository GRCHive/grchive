package core_test

import (
	"encoding/hex"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/mail_api"
	"reflect"
	"strings"
	"testing"
)

func mustDecodeString(str string) []byte {
	b, _ := hex.DecodeString(str)
	return b
}

func TestInitializeConfig(t *testing.T) {
	// Just test that it initializes TomlConfig and EnvConfig to non-null
	// pointers.
	core.InitializeConfig("../../../../src/webserver/config/config.toml")
	assert.NotNil(t, core.TomlConfig)
	assert.NotNil(t, core.EnvConfig)
}

type SingleFieldConfig struct {
	StructFieldName string
	TomlFieldName   string
	Value           interface{}
	ConfigValue     interface{}
}

type GenerateTomlConfig struct {
	Params []SingleFieldConfig
}

func getNestedField(val reflect.Value, field string) reflect.Value {
	fields := strings.Split(field, ".")
	currentVal := val
	for _, f := range fields {
		currentVal = currentVal.FieldByName(f)
		if currentVal.Kind() == reflect.Ptr {
			currentVal = currentVal.Elem()
		}
	}
	return currentVal
}

// Generates a TOML Tree along with the expected EnvConfig.
func generateTestToml(config GenerateTomlConfig) (*toml.Tree, *core.EnvConfigData) {
	newTree, _ := toml.TreeFromMap(map[string]interface{}{})
	newData := core.EnvConfigData{
		Okta:               new(core.OktaConfig),
		Login:              new(core.LoginConfig),
		Company:            new(core.CompanyConfig),
		Vault:              new(core.VaultConfig),
		Gcloud:             new(core.GCloudConfig),
		Mail:               new(core.MailConfig),
		HashId:             new(core.HashIdConfigData),
		RabbitMQ:           new(core.RabbitMQConfig),
		Grpc:               new(core.GrpcEndpoints),
		Tls:                new(core.TLSConfig),
		GitlabRegistryAuth: new(core.ContainerRegistryAuth),
		Kotlin:             new(core.KotlinConfigData),
		Gitea:              new(core.GiteaConfigData),
		Drone:              new(core.DroneConfigData),
		Features:           new(core.FeatureFlags),
		Artifactory:        new(core.ArtifactoryConfigData),
		Notifications:      new(core.NotificationConfig),
	}

	newDataVal := reflect.ValueOf(&newData).Elem()

	for _, paramCfg := range config.Params {
		newTree.Set(paramCfg.TomlFieldName, paramCfg.Value)
		field := getNestedField(newDataVal, paramCfg.StructFieldName)

		if paramCfg.ConfigValue != nil {
			field.Set(reflect.ValueOf(paramCfg.ConfigValue))
		} else {
			field.Set(reflect.ValueOf(paramCfg.Value))
		}
	}

	return newTree, &newData
}

func TestLoadEnvConfig(t *testing.T) {
	for _, ref := range []struct {
		cfg        GenerateTomlConfig
		parseError bool
	}{
		{
			// Empty TOML config should fail.
			cfg: GenerateTomlConfig{
				Params: []SingleFieldConfig{},
			},
			parseError: true,
		},
		{
			// Perfectly constructed TOML with proper
			// types should succeed.
			cfg: GenerateTomlConfig{
				Params: []SingleFieldConfig{
					{"SelfUri", "self_uri", "test_uri", nil},
					{"SelfDomain", "self_domain", "test_domain", nil},
					{"DatabaseConnString", "database.connection", "connection_string", nil},
					{"DatabaseUsername", "database.username", "username", nil},
					{"DatabasePassword", "database.password", "password", nil},
					{"Okta.BaseUrl", "okta.url", "okta_base_url", nil},
					{"Okta.ApiEndpoint", "okta.api_endpoint", "okta_apiendpoint", nil},
					{"Okta.ApiKey", "okta.api_key", "okta_apikey", nil},
					{"Okta.UsersEndpoint", "okta.users_endpoint", "okta_user_endpoint", nil},
					{"Login.AuthServerEndpoint", "login.authserver_endpoint", "authserver_endpoint", nil},
					{"Login.AuthEndpoint", "login.auth_endpoint", "auth_endpoint", nil},
					{"Login.TokenEndpoint", "login.token_endpoint", "token_endpoint", nil},
					{"Login.KeyEndpoint", "login.key_endpoint", "key_endpoint", nil},
					{"Login.LogoutEndpoint", "login.logout_endpoint", "logout_endpoint", nil},
					{"Login.ClientId", "login.params.client_id", "client_id", nil},
					{"Login.ClientSecret", "login.params.client_secret", "client_secret", nil},
					{"Login.ResponseType", "login.params.response_type", "response_type", nil},
					{"Login.ResponseMode", "login.params.response_mode", "response_mode", nil},
					{"Login.Scope", "login.params.scope", "scope", nil},
					{"Login.RedirectUrl", "login.params.redirect_uri", "redirect_uri", nil},
					{"Login.GrantType", "login.params.grant_type", "grant_type", nil},
					{"Login.AuthAudience", "login.auth_audience", "audience", nil},
					{"Login.TimeDriftLeewaySeconds", "login.time_drift_leeway_seconds", int64(60), 60.0},
					{"SessionKeys", "security.session_keys", []interface{}{
						"ABCDEF",
					}, [][]byte{
						mustDecodeString("ABCDEF"),
					}},
					{"HmacKey", "security.hmac_key", "ABCDEF", mustDecodeString("ABCDEF")},
					{"UseSecureCookies", "security.use_secure_cookies", true, nil},
					{"LogEncryptionPath", "security.log_encryption_path", "path", nil},
					{"Company.CompanyName", "company.company_name", "company_name", nil},
					{"Company.Domain", "company.domain", "domain", nil},
					{"Vault.Url", "vault.url", "url", nil},
					{"Vault.Username", "vault.userpass.username", "username", nil},
					{"Vault.Password", "vault.userpass.password", "password", nil},
					{"Gcloud.AuthFilename", "gcloud.credentials_file", "filenamegoeshere", nil},
					{"Gcloud.DocBucket", "gcloud.storage.doc_bucket", "bucketname", nil},
					{"Mail.Provider", "mail.provider", "provider", mail.MailAPIProvider("provider")},
					{"Mail.Key", "mail.key", "key", nil},
					{"Mail.VeriEmailFrom.Name", "mail.verification.from.name", "name", nil},
					{"Mail.VeriEmailFrom.Email", "mail.verification.from.email", "email", nil},
					{"HashId.Salt", "hashids.salt", "salt", nil},
					{"HashId.MinLength", "hashids.min_length", int64(100), 100},
					{"RabbitMQ.Username", "rabbitmq.username", "name", nil},
					{"RabbitMQ.Password", "rabbitmq.password", "asdfasdf", nil},
					{"RabbitMQ.Host", "rabbitmq.host", "hostname", nil},
					{"RabbitMQ.Port", "rabbitmq.port", int64(64), int32(64)},
					{"RabbitMQ.UseTLS", "rabbitmq.use_tls", true, nil},
					{"Grpc.QueryRunner.Host", "grpc.query_runner.host", "hostname", nil},
					{"Grpc.QueryRunner.Port", "grpc.query_runner.port", int64(128), int32(128)},
					{"Grpc.QueryRunner.TLS", "grpc.query_runner.tls.enable", true, nil},
					{"Grpc.QueryRunner.TLSCert", "grpc.query_runner.tls.cert", "certfile", nil},
					{"Grpc.QueryRunner.TLSKey", "grpc.query_runner.tls.key", "keyfile", nil},
					{"Tls.TLSRootCaCert", "tls.root_ca", "root_ca", nil},
					{"GitlabRegistryAuth.Username", "registry.gitlab.username", "testuser", nil},
					{"GitlabRegistryAuth.Password", "registry.gitlab.password", "testpassword", nil},
					{"Kotlin.GroupId", "kotlin.group_id", "testbucket", nil},
					{"Kotlin.ArtifactId", "kotlin.artifact_id", "testbucket", nil},
					{"Kotlin.MajorVersion", "kotlin.major_version", int64(10), int(10)},
					{"Kotlin.MinorVersion", "kotlin.minor_version", int64(20), int(20)},
					{"Gitea.Host", "gitea.host", "host", nil},
					{"Gitea.Port", "gitea.port", int64(128), int32(128)},
					{"Gitea.Protocol", "gitea.protocol", "protocl", nil},
					{"Gitea.Token", "gitea.token", "token", nil},
					{"Gitea.GlobalOrg", "gitea.global_org", "grchive", nil},
					{"Features.Automation", "features.automation", true, nil},
					{"Drone.Host", "drone.host", "host", nil},
					{"Drone.Port", "drone.port", int64(128), int32(128)},
					{"Drone.Protocol", "drone.protocol", "protocl", nil},
					{"Drone.Token", "drone.token", "token", nil},
					{"Drone.RunnerType", "drone.runner_type", "typetype", nil},
					{"Drone.RunnerImage", "drone.runner_image", "typetype", nil},
					{"Drone.RunnerImagePull", "drone.runner_image_pull", "typetype", nil},
					{"ScriptRunner.RunnerImage", "script_runner.runner_image", "typetype", nil},
					{"Artifactory.Host", "artifactory.host", "host", nil},
					{"Artifactory.Port", "artifactory.port", int64(128), int32(128)},
					{"Notifications.EnableEmail", "notifications.enable_email", true, nil},
				},
			},
			parseError: false,
		},
	} {
		refTree, refConfig := generateTestToml(ref.cfg)

		if ref.parseError {
			assert.Panics(t, func() { core.LoadEnvConfig(refTree) })
		} else {
			assert.True(t, reflect.DeepEqual(
				refConfig,
				core.LoadEnvConfig(refTree)))
		}
	}
}
