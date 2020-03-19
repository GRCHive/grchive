package main

import (
	"encoding/json"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/mail_api"
	"gitlab.com/grchive/grchive/webcore"
)

func handleGettingStartedEvent(event *core.Event) *webcore.RabbitMQError {
	// Special case where we just want to send an email to relevant people
	// at GRCHive...don't particularly care about storing this in the DB.
	mailPayload := mail.MailPayload{
		From: core.EnvConfig.Mail.VeriEmailFrom,
		To: mail.Email{
			Name:  "GRCHive Sales",
			Email: "sales@grchive.com",
		},
		Subject: "GRCHive - New Getting Started Entry",
		Message: fmt.Sprintf("%s (%s) filled out the Getting Started form.", event.Subject.FirstName, event.Subject.Email),
	}

	err := mail.MailProvider.SendMail(mailPayload)
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}
	return nil
}

func generateNotification(data []byte) *webcore.RabbitMQError {
	core.Info(string(data))

	incomingMessage := webcore.EventMessage{}
	err := json.Unmarshal(data, &incomingMessage)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}

	event, err := incomingMessage.RecreateEvent()
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}

	if event.Verb == core.VerbGettingStarted {
		return handleGettingStartedEvent(event)
	}

	core.Debug("\tFind Relevant Users")
	relevantUsers, err := webcore.FindRelevantUsersForEvent(event)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}

	// If we don't need to notify anyone then don't do any further work.
	if len(relevantUsers) == 0 {
		return nil
	}

	notification, err := webcore.CreateNotificationFromEvent(event)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}

	tx := database.CreateTx()

	core.Debug("\tInsert Notification")
	err = database.InsertNotificationWithTx(notification, tx)
	if err != nil {
		tx.Rollback()
		return &webcore.RabbitMQError{err, false}
	}

	core.Debug("\tLink Notification")
	err = database.LinkNotificationToUsersWithTx(notification.Id, notification.OrgId, relevantUsers, tx)
	if err != nil {
		tx.Rollback()
		return &webcore.RabbitMQError{err, false}
	}

	core.Debug("\tCommit")
	err = tx.Commit()
	if err != nil {
		return &webcore.RabbitMQError{err, true}
	}

	webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
		Exchange: webcore.NOTIFICATION_EXCHANGE,
		Queue:    "",
		Body: webcore.NotificationMessage{
			Notification:  *notification,
			RelevantUsers: relevantUsers,
		},
	})

	return nil
}

func main() {
	core.Init()
	database.Init()
	mail.InitializeMailAPI(core.EnvConfig.Mail.Provider, core.EnvConfig.Mail.Key)

	webcore.DefaultRabbitMQ.Connect(*core.EnvConfig.RabbitMQ, webcore.QueueConfig{}, core.EnvConfig.Tls)
	defer webcore.DefaultRabbitMQ.Cleanup()

	forever := make(chan bool)

	webcore.DefaultRabbitMQ.ReceiveMessages(webcore.EVENT_NOTIFICATION_QUEUE, generateNotification)

	<-forever
}
