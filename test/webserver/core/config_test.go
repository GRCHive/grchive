package core_test

import (
	"encoding/hex"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"gitlab.com/b3h47pte/audit-stuff/core"
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
	core.InitializeConfig("../../../src/webserver/config/config.toml")
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
		Okta:    new(core.OktaConfig),
		Login:   new(core.LoginConfig),
		Company: new(core.CompanyConfig),
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
					{"DatabaseConnString", "database.connection", "connection_string", nil},
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
					{"UseSecureCookies", "security.use_secure_cookies", true, nil},
					{"Company.CompanyName", "company.company_name", "company_name", nil},
					{"Company.Domain", "company.domain", "domain", nil},
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
