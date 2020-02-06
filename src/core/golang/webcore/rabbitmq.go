package webcore

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/grchive/grchive/core"
)

const DEFAULT_EXCHANGE string = ""
const FILE_PREVIEW_QUEUE string = "filepreview"
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

type RecvMsgFn func([]byte) *RabbitMQError

func SetupChannel(channel *amqp.Channel) {
	_, err := channel.QueueDeclare(
		FILE_PREVIEW_QUEUE, // name
		true,               // durable
		false,              // auto delete
		false,              // exclusive
		false,              // no wait
		nil,                // arguments
	)

	if err != nil {
		core.Error("Failed to declare queue: " + err.Error())
	}
}

type PublishMessage struct {
	Exchange string
	Queue    string
	Body     interface{}
}

type RabbitMQInterface interface {
	// Setup functions
	Connect(core.RabbitMQConfig)
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

func CreateRabbitMQConnection(url string, numChannels int) *RabbitMQConnection {
	connection, err := amqp.Dial(url)
	if err != nil {
		core.Error("Failed to dial RabbitMQ: " + err.Error())
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
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
}

func (r *RealRabbitMQInterface) Connect(cfg core.RabbitMQConfig) {
	url := generateUrlFromConfig(cfg)
	r.publish = CreateRabbitMQConnection(url, 4)
	r.consume = CreateRabbitMQConnection(url, 1)
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
