package vault

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type VaultConfig struct {
	Url      string
	Username string
	Password string

	token  string
	client *http.Client
}

var GlobalConfig VaultConfig = VaultConfig{}

func Initialize(config VaultConfig, tlsConfig *tls.Config) {
	GlobalConfig = config

	GlobalConfig.client = &http.Client{}
	if tlsConfig != nil {
		GlobalConfig.client.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	if err := userPassAuth(config.Username, config.Password); err != nil {
		panic("Failed to auth with Vault: " + err.Error())
	}
}

func sendVaultRequest(method string, endpoint string, data interface{}) (map[string]*json.RawMessage, error) {
	body := &bytes.Buffer{}
	if data != nil {
		jsonBuffer, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(jsonBuffer)
	}

	fullUrl := GlobalConfig.Url + endpoint
	req, err := http.NewRequest(method, fullUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Vault-Token", GlobalConfig.token)
	resp, err := GlobalConfig.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		rootObj := map[string]*json.RawMessage{}
		err = json.Unmarshal(respBodyData, &rootObj)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Vault Response: %s [Parse Fail: %s]", string(respBodyData), err.Error()))
		}

		rawErrors, ok := rootObj["errors"]
		auxString := ""
		if ok {
			errors := make([]string, 0)
			err = json.Unmarshal(*rawErrors, &errors)
			if err == nil {
				for _, s := range errors {
					auxString += s + " :: "
				}
			}
		}
		return nil, errors.New(fmt.Sprintf("%s -> %d -- Vault Request Failed: %s", fullUrl, resp.StatusCode, auxString))
	} else {
		rootObj := map[string]*json.RawMessage{}
		// Don't handle errors error since we might not actually have a body to parse.
		json.Unmarshal(respBodyData, &rootObj)
		return rootObj, nil
	}
}
