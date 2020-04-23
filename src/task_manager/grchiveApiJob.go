package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type GrchiveApiJobHandler struct {
	taskId int64
	data   interface{}
	userId int64
}

func sendGrchiveApiRequest(endpoint string, method string, body interface{}, apiKey core.RawApiKey) error {
	client := http.Client{}

	sendBody := &bytes.Buffer{}
	if body != nil {
		jsonBuffer, err := json.Marshal(body)
		if err != nil {
			return err
		}
		sendBody = bytes.NewBuffer(jsonBuffer)
	}

	// Do we want to route this via our own internal hostname instead?
	fullUrl := core.EnvConfig.SelfUri + endpoint

	req, err := http.NewRequest(method, fullUrl, sendBody)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ApiKey", string(apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Failed to send API request.")
	}

	return nil
}

func (h *GrchiveApiJobHandler) Tick(c core.Clock) error {
	rawData, err := json.Marshal(h.data)
	if err != nil {
		return err
	}

	data := core.GrchiveApiTaskData{}
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return err
	}

	// All API tasks will have "ScheduledTaskId" added to its payload
	// so that the API endpoint can identify which scheduled task
	// is the reason it was hit (as opposed to manually be a user).
	rawPayload, err := json.Marshal(data.Payload)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{}
	err = json.Unmarshal(rawPayload, &payload)
	if err != nil {
		return err
	}
	payload["ScheduledTaskId"] = h.taskId

	// Create a temporary API key for the user to use just for this API request.
	rawKey, apiKey := webcore.GenerateTemporaryAPIKeyForUser(h.userId)
	err = database.StoreApiKey(apiKey)
	if err != nil {
		return err
	}

	err = sendGrchiveApiRequest(data.Endpoint, data.Method, payload, rawKey)
	if err != nil {
		return err
	}

	return database.DeleteApiKey(apiKey.HashedKey)
}
