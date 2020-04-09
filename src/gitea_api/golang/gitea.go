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
) (int, map[string]*json.RawMessage, error) {
	body := &bytes.Buffer{}
	if data != nil {
		jsonBuffer, err := json.Marshal(data)
		if err != nil {
			return -1, nil, err
		}

		body = bytes.NewBuffer(jsonBuffer)
	}

	req, err := http.NewRequest(method, fullUrl, body)
	if err != nil {
		return -1, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()

	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, nil, err
	}

	rootObj := map[string]*json.RawMessage{}
	err = json.Unmarshal(respBodyData, &rootObj)

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusAccepted &&
		resp.StatusCode != http.StatusNoContent {
		// The Unmarshal might fail if there's no body.
		// That's up to whoever calls this function to handle/know based off the API spec.
		return resp.StatusCode, rootObj, errors.New(fmt.Sprintf("%s -> %d -- Gitea Request Failed: %s", fullUrl, resp.StatusCode, string(respBodyData)))
	}

	return resp.StatusCode, rootObj, err
}

func (r RealGiteaApi) sendGiteaRequestWithUserAuth(
	method string,
	endpoint string,
	user GiteaUser,
	data interface{},
) (int, map[string]*json.RawMessage, error) {
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
) (int, map[string]*json.RawMessage, error) {
	return r.sendGiteaRequest(
		method,
		r.cfg.apiUrl()+endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("token %s", token),
		},
		data,
	)
}
