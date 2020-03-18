package main

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
)

func generateNotification(data []byte) *webcore.RabbitMQError {
	core.Info(string(data))

	incomingMessage := webcore.EventMessage{}
	err := json.Unmarshal(data, &incomingMessage)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}

	webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
		Exchange: webcore.NOTIFICATION_EXCHANGE,
		Queue:    "",
		Body: webcore.NotificationMessage{
			Notification: core.Notification{},
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
