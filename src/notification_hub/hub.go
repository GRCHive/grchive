package main

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
)

func handleGettingStartedEvent(event *core.Event) *webcore.RabbitMQError {
	// Special case where we just want to send an email to relevant people
	// at GRCHive...don't particularly care about storing this in the DB.
	return nil
}

func generateNotification(data []byte) *webcore.RabbitMQError {
	core.Info(string(data))

	incomingMessage := webcore.EventMessage{}
	err := json.Unmarshal(data, &incomingMessage)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}

	if incomingMessage.Event.Verb == core.VerbGettingStarted {
		return handleGettingStartedEvent(&incomingMessage.Event)
	}

	notification, err := webcore.CreateNotificationFromEvent(&incomingMessage.Event)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}

	relevantUsers, err := webcore.FindRelevantUsersForEvent(&incomingMessage.Event)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}

	tx := database.CreateTx()

	err = database.InsertNotificationWithTx(notification, tx)
	if err != nil {
		tx.Rollback()
		return &webcore.RabbitMQError{err, false}
	}

	err = database.LinkNotificationToUsersWithTx(notification.Id, notification.OrgId, relevantUsers, tx)
	if err != nil {
		tx.Rollback()
		return &webcore.RabbitMQError{err, false}
	}

	err = tx.Commit()
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}

	webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
		Exchange: webcore.NOTIFICATION_EXCHANGE,
		Queue:    "",
		Body: webcore.NotificationMessage{
			Notification: *notification,
		},
	})

	return nil
}

func main() {
	core.Init()
	database.Init()

	webcore.DefaultRabbitMQ.Connect(*core.EnvConfig.RabbitMQ, core.EnvConfig.Tls)
	defer webcore.DefaultRabbitMQ.Cleanup()

	forever := make(chan bool)

	webcore.DefaultRabbitMQ.ReceiveMessages(webcore.EVENT_NOTIFICATION_QUEUE, generateNotification)

	<-forever
}
