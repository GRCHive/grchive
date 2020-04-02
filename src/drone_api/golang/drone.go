package drone

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RealDroneApi struct {
	cfg    DroneConfig
	client *http.Client
}

var GlobalDroneApi = RealDroneApi{}

func (d *RealDroneApi) MustInitialize(cfg DroneConfig) {
	d.cfg = cfg
	d.client = &http.Client{}
}

func (r RealDroneApi) sendDroneRequest(
	method string,
	endpoint string,
	token string,
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

	fullUrl := r.cfg.apiUrl() + endpoint
	req, err := http.NewRequest(method, fullUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusAccepted &&
		resp.StatusCode != http.StatusNoContent {
		return nil, errors.New(fmt.Sprintf("%s -> %d -- Drone Request Failed: %s", fullUrl, resp.StatusCode, string(respBodyData)))
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
