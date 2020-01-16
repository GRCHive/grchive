package webcore

import (
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

type RabbitMQInterface interface {
	Connect(cfg core.RabbitMQConfig)
}

type RealRabbitMQInterface struct {
	connection *amqp.Connection
}

func generateUrlFromConfig(cfg core.RabbitMQConfig) string {
	return fmt.Sprintf("amqp://%s:%d/", cfg.Host, cfg.Port)
}

func (r *RealRabbitMQInterface) Connect(cfg core.RabbitMQConfig) {
	var err error
	r.connection, err = amqp.Dial(generateUrlFromConfig(cfg))
	if err != nil {
		panic("Failed to connect to RabbitMQ: " + err.Error())
	}
}

var DefaultRabbitMQ = RealRabbitMQInterface{}
