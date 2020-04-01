package gitea

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RealGiteaApi struct {
	cfg    GiteaConfig
	client *http.Client
}

var GlobalGiteaApi = RealGiteaApi{}

func (r RealGiteaApi) sendGiteaRequest(
	method string,
	fullUrl string,
	headers map[string]string,
	data interface{},
) (map[string]*json.RawMessage, error) {
	body := &bytes.Buffer{}
	if data != nil {
		jsonBuffer, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(jsonBuffer)
	}

	req, err := http.NewRequest(method, fullUrl, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return nil, errors.New(fmt.Sprintf("%s -> %d -- Gitea Request Failed: %s", fullUrl, resp.StatusCode, string(respBodyData)))
	}

	rootObj := map[string]*json.RawMessage{}
	err = json.Unmarshal(respBodyData, &rootObj)
	if err != nil {
		// This isn't necessarily an error since not all
		// requests will return something.
		return rootObj, nil
	}

	return rootObj, nil
}

func (r RealGiteaApi) sendGiteaRequestWithUserAuth(
	method string,
	endpoint string,
	user GiteaUser,
	data interface{},
) (map[string]*json.RawMessage, error) {
	return r.sendGiteaRequest(
		method,
		r.cfg.apiUrlUserAuth(user)+endpoint,
		map[string]string{},
		data,
	)
}

func (r RealGiteaApi) sendGiteaRequestWithToken(
	method string,
	endpoint string,
	token string,
	data interface{},
) (map[string]*json.RawMessage, error) {
	return r.sendGiteaRequest(
		method,
		r.cfg.apiUrl()+endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("token %s", token),
		},
		data,
	)

}
