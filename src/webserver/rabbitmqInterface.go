package main

import (
	"encoding/json"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/mail_api"
	"gitlab.com/grchive/grchive/webcore"
	"html/template"
	"strconv"
	"time"
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

	sendNotification(&incomingMessage.Notification, incomingMessage.RelevantUsers)

	return nil
}

func sendNotification(notification *core.Notification, users []*core.User) error {
	emailTemplate, err := template.ParseFiles("src/webserver/templates/email/notification.tmpl")
	if err != nil {
		return err
	}

	subjectHandle, err := webcore.GetResourceHandle(notification.SubjectType, notification.SubjectId, notification.OrgId)
	if err != nil {
		return err
	}

	objectHandle, err := webcore.GetResourceHandle(notification.ObjectType, notification.ObjectId, notification.OrgId)
	if err != nil {
		return err
	}

	indirectObjectHandle, err := webcore.GetResourceHandle(notification.IndirectObjectType, notification.IndirectObjectId, notification.OrgId)
	if err != nil {
		return err
	}

	org, err := database.FindOrganizationFromId(notification.OrgId)
	if err != nil {
		return err
	}

	message, err := core.HtmlTemplateToString(emailTemplate, map[string]string{
		"subject":        subjectHandle.DisplayText,
		"verb":           notification.Verb,
		"object":         objectHandle.DisplayText,
		"objectUrl":      objectHandle.ResourceUri.NullString.String,
		"indirectObject": indirectObjectHandle.DisplayText,
		"timestamp":      notification.Time.Format(time.UnixDate),
	})

	if err != nil {
		return err
	}

	for _, u := range users {
		mailPayload := mail.MailPayload{
			From: core.EnvConfig.Mail.VeriEmailFrom,
			To: mail.Email{
				Name:  fmt.Sprintf("%s %s", u.FirstName, u.LastName),
				Email: u.Email,
			},
			Subject: fmt.Sprintf("Notification in %s for %s", org.Name, objectHandle.DisplayText),
			Message: message,
		}

		mail.MailProvider.SendMail(mailPayload)
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
