package vault

import (
	"encoding/json"
	"errors"
	"strings"
)

const SecretDataEndpoint string = "/v1/secret/data/"

func StoreSecret(path string, data map[string]string, cas int) error {
	options := map[string]interface{}{}
	if cas != -1 {
		options["cas"] = cas
	}

	params := map[string]interface{}{
		"options": options,
		"data":    data,
	}

	_, err := sendVaultRequest("POST", SecretDataEndpoint+strings.TrimPrefix(path, "secret/"), params)
	return err
}

func GetSecret(path string) (map[string]string, error) {
	data, err := sendVaultRequest("GET", SecretDataEndpoint+strings.TrimPrefix(path, "secret/"), nil)
	if err != nil {
		return nil, err
	}

	// The secret data is stored in ['data']['data'].
	dataMap := map[string]*json.RawMessage{}
	err = json.Unmarshal(*data["data"], &dataMap)
	if err != nil {
		return nil, err
	}

	secretMap := map[string]string{}
	err = json.Unmarshal(*dataMap["data"], &secretMap)
	if err != nil {
		return nil, err
	}

	return secretMap, nil
}

func GetSecretWithKey(path string, key string) (string, error) {
	data, err := GetSecret(path)
	if err != nil {
		return "", err
	}

	ret, ok := data[key]
	if !ok {
		return "", errors.New("Failed to find key in secret.")
	}

	return ret, nil
}
