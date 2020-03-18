package main

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"time"
)

func onNotifyControlOwnerChange(data string) error {
	type Data struct {
		User    core.User
		Control core.Control
	}

	parsedData := Data{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return err
	}

	if !parsedData.Control.OwnerId.NullInt64.Valid {
		return nil
	}

	assignedToUser, err := database.FindUserFromId(parsedData.Control.OwnerId.NullInt64.Int64)
	if err != nil {
		return err
	}

	webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
		Exchange: webcore.EVENT_EXCHANGE,
		Queue:    "",
		Body: webcore.EventMessage{
			Event: core.Event{
				Subject:        parsedData.User,
				Verb:           core.VerbAssign,
				Object:         assignedToUser,
				IndirectObject: parsedData.Control,
				Timestamp:      time.Now().UTC(),
			},
		},
	})

	return nil
}
