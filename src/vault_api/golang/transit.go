package vault

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

const TransitCreateKeyEndpoint string = "/v1/transit/keys/"
const TransitEncryptEndpoint string = "/v1/transit/encrypt/"
const TransitDecryptEndpoint string = "/v1/transit/decrypt/"

func TransitCreateNewEngineKey(path string) error {
	_, err := sendVaultRequest("POST", TransitCreateKeyEndpoint+path, nil)
	if err != nil {
		return err
	}
	return nil
}

func TransitEncrypt(path string, data []byte) ([]byte, error) {
	reqData := map[string]string{
		"plaintext": base64.StdEncoding.EncodeToString(data),
	}

	respData, err := sendVaultRequest("POST", TransitEncryptEndpoint+path, reqData)
	if err != nil {
		return nil, err
	}

	dataBlock, ok := respData["data"]
	if !ok {
		return nil, errors.New("No data in encryption response")
	}

	type ParsedData struct {
		Ciphertext string `json:"ciphertext"`
	}
	parsed := ParsedData{}
	err = json.Unmarshal(*dataBlock, &parsed)
	if err != nil {
		return nil, err
	}

	return []byte(parsed.Ciphertext), nil
}

func TransitDecrypt(path string, data []byte) ([]byte, error) {
	reqData := map[string]string{
		"ciphertext": string(data),
	}

	respData, err := sendVaultRequest("POST", TransitDecryptEndpoint+path, reqData)
	if err != nil {
		return nil, err
	}

	dataBlock, ok := respData["data"]
	if !ok {
		return nil, errors.New("No data in encryption response")
	}

	type ParsedData struct {
		Plaintext string `json:"plaintext"`
	}
	parsed := ParsedData{}
	err = json.Unmarshal(*dataBlock, &parsed)
	if err != nil {
		return nil, err
	}

	ret, err := base64.StdEncoding.DecodeString(parsed.Plaintext)
	return ret, err
}
