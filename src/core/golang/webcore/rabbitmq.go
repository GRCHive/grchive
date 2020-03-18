package webcore

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/grchive/grchive/core"
)

const DEFAULT_EXCHANGE string = ""
const EVENT_EXCHANGE string = "events"
const NOTIFICATION_EXCHANGE string = "notifications"

const FILE_PREVIEW_QUEUE string = "filepreview"
const DATABASE_REFRESH_QUEUE string = "dbrefresh"
const EVENT_NOTIFICATION_QUEUE string = "eventnotification"

const CHANNEL_BUFFER int = 12

type RabbitMQError struct {
	Err     error
	Requeue bool
}

func (e RabbitMQError) Error() string {
	if e.Err == nil {
		return "No Error"
	}
	return e.Err.Error()
}

// Pre-defined message types.
type FilePreviewMessage struct {
	File    core.ControlDocumentationFile
	Storage core.FileStorageData
}

type DatabaseRefreshMessage struct {
	Refresh core.DbRefresh
}

type EventMessage struct {
	Event              core.Event
	ObjectType         string
	IndirectObjectType string
}

type NotificationMessage struct {
	Notification core.Notification
}

type RecvMsgFn func([]byte) *RabbitMQError

func SetupChannel(channel *amqp.Channel) {
	//
	// Default Exchange
	//
	_, err := channel.QueueDeclare(
		FILE_PREVIEW_QUEUE, // name
		true,               // durable
		false,              // auto delete
		false,              // exclusive
		false,              // no wait
		nil,                // arguments
	)

	if err != nil {
		core.Error("Failed to declare file preview queue: " + err.Error())
	}

	_, err = channel.QueueDeclare(
		DATABASE_REFRESH_QUEUE, // name
		true,                   // durable
		false,                  // auto delete
		false,                  // exclusive
		false,                  // no wait
		nil,                    // arguments
	)

	if err != nil {
		core.Error("Failed to declare database refresh queue: " + err.Error())
	}

	//
	// Event Exchange
	//

	err = channel.ExchangeDeclare(
		EVENT_EXCHANGE, // name
		"fanout",       // type
		true,           // durable
		false,          // auto delete
		false,          // internal
		false,          // no wait
		nil,            // arguments
	)

	if err != nil {
		core.Error("Failed to declare event exchange: " + err.Error())
	}

	_, err = channel.QueueDeclare(
		EVENT_NOTIFICATION_QUEUE, // name
		true,                     // durable
		false,                    // auto delete
		false,                    // exclusive
		false,                    // no wait
		nil,                      // arguments
	)

	if err != nil {
		core.Error("Failed to declare event notification queue: " + err.Error())
	}

	err = channel.QueueBind(
		EVENT_NOTIFICATION_QUEUE, // queue name
		"",                       // routing key
		EVENT_EXCHANGE,           // exchange
		false,                    // no wait
		nil,                      // arguments
	)

	if err != nil {
		core.Error("Failed to bind event notification queue to exchange: " + err.Error())
	}

	//
	// Notification Exchange
	//

	err = channel.ExchangeDeclare(
		NOTIFICATION_EXCHANGE, // name
		"fanout",              // type
		false,                 // durable
		false,                 // auto delete
		false,                 // internal
		false,                 // no wait
		nil,                   // arguments
	)

	if err != nil {
		core.Error("Failed to declare notification exchange: " + err.Error())
	}

	q, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no wait
		nil,   // arguments
	)

	if err != nil {
		core.Error("Failed to declare notification queue: " + err.Error())
	}

	err = channel.QueueBind(
		q.Name,                // queue name
		"",                    // routing key
		NOTIFICATION_EXCHANGE, // exchange
		false,                 // no wait
		nil,                   // arguments
	)

	if err != nil {
		core.Error("Failed to bind notification queue to exchange: " + err.Error())
	}
}

type PublishMessage struct {
	Exchange string
	Queue    string
	Body     interface{}
}

type RabbitMQInterface interface {
	// Setup functions
	Connect(core.RabbitMQConfig, *core.TLSConfig)
	Cleanup()

	// Message IO
	SendMessage(PublishMessage)
	ReceiveMessages(string, RecvMsgFn) error
}

type RabbitMQConnection struct {
	Connection *amqp.Connection
	Channels   []*amqp.Channel

	publishChannel chan PublishMessage
}

func (r *RabbitMQConnection) publishWorker(idx int) {
	for {
		msg := <-r.publishChannel

		byteMsg, err := json.Marshal(msg.Body)
		if err != nil {
			core.Error("Failed to marshal message: " + err.Error())
		}

		err = r.Channels[idx].Publish(
			msg.Exchange,
			msg.Queue,
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         byteMsg,
				DeliveryMode: amqp.Persistent,
			},
		)

		if err != nil {
			core.Warning("Failed to publish : " + err.Error())
		}
	}
}

func CreateRabbitMQConnection(cfg core.RabbitMQConfig, tls *core.TLSConfig, numChannels int) *RabbitMQConnection {
	var connection *amqp.Connection
	var err error
	url := generateUrlFromConfig(cfg)

	if cfg.UseTLS && tls != nil {
		connection, err = amqp.DialTLS(url, tls.Config())
		if err != nil {
			core.Error("Failed to dial RabbitMQ (TLS): " + err.Error())
		}
	} else {
		connection, err = amqp.Dial(url)
		if err != nil {
			core.Error("Failed to dial RabbitMQ (non-TLS): " + err.Error())
		}
	}

	c := RabbitMQConnection{
		Connection:     connection,
		Channels:       make([]*amqp.Channel, numChannels),
		publishChannel: make(chan PublishMessage, CHANNEL_BUFFER),
	}

	for i := 0; i < numChannels; i++ {
		c.Channels[i], err = connection.Channel()
		if err != nil {
			core.Error("Failed to create channel: " + err.Error())
		}
		SetupChannel(c.Channels[i])
		go c.publishWorker(i)
	}

	return &c
}

func (c *RabbitMQConnection) SendMessage(msg PublishMessage) {
	c.publishChannel <- msg
}

func (c *RabbitMQConnection) ReceiveMessages(queue string, fn RecvMsgFn) error {
	msgs, err := c.Channels[0].Consume(
		queue,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		var err error
		for d := range msgs {
			mqErr := fn(d.Body)
			if mqErr != nil {
				core.Warning("Failed to process: " + mqErr.Error())
				err = d.Nack(false, mqErr.Requeue)
			} else {
				err = d.Ack(false)
			}

			if err != nil {
				core.Error("Failed to ack: " + err.Error())
			}
		}
	}()

	return nil
}

func (c *RabbitMQConnection) Close() {
	for i := 0; i < len(c.Channels); i++ {
		c.Channels[i].Close()
	}
	c.Connection.Close()
}

type RealRabbitMQInterface struct {
	publish *RabbitMQConnection
	consume *RabbitMQConnection
}

func generateUrlFromConfig(cfg core.RabbitMQConfig) string {
	var prefix string
	if cfg.UseTLS {
		prefix = "amqps://"
	} else {
		prefix = "amqp://"
	}
	return fmt.Sprintf("%s%s:%s@%s:%d/", prefix, cfg.Username, cfg.Password, cfg.Host, cfg.Port)
}

func (r *RealRabbitMQInterface) Connect(cfg core.RabbitMQConfig, tls *core.TLSConfig) {
	r.publish = CreateRabbitMQConnection(cfg, tls, 4)
	r.consume = CreateRabbitMQConnection(cfg, tls, 1)
}

func (r *RealRabbitMQInterface) SendMessage(msg PublishMessage) {
	r.publish.SendMessage(msg)
}

func (r *RealRabbitMQInterface) ReceiveMessages(queue string, fn RecvMsgFn) error {
	return r.consume.ReceiveMessages(queue, fn)
}

func (r *RealRabbitMQInterface) Cleanup() {
	r.publish.Close()
	r.consume.Close()
}

var DefaultRabbitMQ = RealRabbitMQInterface{}
