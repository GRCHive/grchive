package security

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"io/ioutil"
	"net/http"
)

const TransitCreateKeyEndpoint string = "/v1/transit/keys/"
const TransitEncryptEndpoint string = "/v1/transit/encrypt/"

func sendVaultRequest(method string, endpoint string, data interface{}) (map[string]*json.RawMessage, error) {
	body := &bytes.Buffer{}
	if data != nil {
		jsonBuffer, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(jsonBuffer)
	}

	fullUrl := core.EnvConfig.Vault.Url + endpoint
	req, err := http.NewRequest(method, fullUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Vault-Token", core.EnvConfig.Vault.Token)
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
	return nil, nil
}
