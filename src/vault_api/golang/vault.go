package vault

import (
	"bytes"
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
	token    string
}

var GlobalConfig VaultConfig = VaultConfig{}

func Initialize(config VaultConfig) {
	GlobalConfig = config
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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rootObj := map[string]*json.RawMessage{}
	err = json.Unmarshal(respBodyData, &rootObj)
	if err != nil {
		// This isn't necessarily an error since not all
		// requests will return something.
		return rootObj, nil
	}

	if resp.StatusCode != http.StatusOK {
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
	}
	return rootObj, nil
}
