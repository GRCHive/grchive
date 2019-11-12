package okta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OktaConfig struct {
	ApiKey    string
	ApiDomain string
}

var GlobalConfig OktaConfig = OktaConfig{}

func InitializeOktaAPI(config OktaConfig) {
	GlobalConfig = config
}

func constructFullUrl(endpoint string) string {
	return GlobalConfig.ApiDomain + endpoint
}

func oktaPost(endpoint string, data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	jsonBuffer := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest("POST", constructFullUrl(endpoint), jsonBuffer)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", GlobalConfig.ApiKey))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
