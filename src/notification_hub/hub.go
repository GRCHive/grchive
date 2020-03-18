package main

import (
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
)

func generateNotification(data []byte) *webcore.RabbitMQError {
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
