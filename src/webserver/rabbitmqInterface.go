package main

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"strconv"
)

func receiveNotification(data []byte) *webcore.RabbitMQError {
	core.Info(string(data))

	incomingMessage := webcore.NotificationMessage{}
	err := json.Unmarshal(data, &incomingMessage)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}

	org, err := database.FindOrganizationFromId(incomingMessage.Notification.OrgId)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}

	// This notification should be fresh so we can just assume that it hasn't been read.
	for _, u := range incomingMessage.RelevantUsers {
		core.DefaultMessageHub.SendMessage(
			core.MHUserNotification,
			core.MessageSubtype(strconv.FormatInt(u.Id, 10)),
			core.NotificationWrapper{
				Notification: incomingMessage.Notification,
				OrgName:      org.OktaGroupName,
				Read:         false,
			},
		)
	}

	return nil
}

// Anything that we want to read from RabbitMQ
// should be handled in this file and send to the message hub.
func SetupRabbitMQInterface() {
	webcore.DefaultRabbitMQ.ReceiveMessages(
		webcore.DefaultRabbitMQ.GetConsumerQueueName(webcore.NotificationQueueId),
		receiveNotification)
}
