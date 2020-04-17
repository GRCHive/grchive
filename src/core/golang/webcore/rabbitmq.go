package webcore

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/grchive/grchive/core"
	"sync/atomic"
	"time"
)

const DEFAULT_EXCHANGE string = ""
const EVENT_EXCHANGE string = "events"
const NOTIFICATION_EXCHANGE string = "notifications"

const FILE_PREVIEW_QUEUE string = "filepreview"
const DATABASE_REFRESH_QUEUE string = "dbrefresh"
const EVENT_NOTIFICATION_QUEUE string = "eventnotification"
const SCRIPT_RUNNER_QUEUE string = "scriptrun"

const (
	NotificationQueueId int = iota
)

const CHANNEL_BUFFER int = 12

type RabbitMQError struct {
	Err     error
	Requeue bool
}

type AmqpChannelWrapper struct {
	Channel *amqp.Channel
	Queues  map[int]string
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
	Notification  core.Notification
	RelevantUsers []*core.User
}

type ScriptRunnerMessage struct {
	RunId int64
	Jar   string
}

type RecvMsgFn func([]byte) *RabbitMQError

func SetupChannel(channel *amqp.Channel, cfg MQClientConfig, idx int, isConsume bool) *AmqpChannelWrapper {
	wrapper := AmqpChannelWrapper{
		Channel: channel,
		Queues:  map[int]string{},
	}

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

	_, err = channel.QueueDeclare(
		SCRIPT_RUNNER_QUEUE, // name
		true,                // durable
		false,               // auto delete
		false,               // exclusive
		false,               // no wait
		nil,                 // arguments
	)

	if err != nil {
		core.Error("Failed to declare script runner queue: " + err.Error())
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

	if cfg.NotificationConsume && idx == 0 && isConsume {
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

		wrapper.Queues[NotificationQueueId] = q.Name
	}

	return &wrapper
}

type PublishMessage struct {
	Exchange string
	Queue    string
	Body     interface{}
}

type MQClientConfig struct {
	// Queue Enabled Flags
	NotificationConsume bool

	// QoS
	ConsumerQos int
}

type RabbitMQInterface interface {
	// Setup functions
	Connect(core.RabbitMQConfig, MQClientConfig, *core.TLSConfig)
	Cleanup()

	// Message IO
	SendMessage(PublishMessage)
	ReceiveMessages(string, RecvMsgFn) error
	WaitForAllMessagesToBeSent()

	// Consumer Queue Names
	GetConsumerQueueName(id int) string
}

type RabbitMQConnection struct {
	Connection *amqp.Connection
	Channels   []*AmqpChannelWrapper

	publishChannel    chan PublishMessage
	publishInProgress int64
}

func (r *RabbitMQConnection) publishWorker(idx int) {
	for {
		msg := <-r.publishChannel

		byteMsg, err := json.Marshal(msg.Body)
		if err != nil {
			atomic.AddInt64(&r.publishInProgress, -1)
			core.Error("Failed to marshal message: " + err.Error())
		}

		err = r.Channels[idx].Channel.Publish(
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
		atomic.AddInt64(&r.publishInProgress, -1)

		if err != nil {
			core.Warning("Failed to publish : " + err.Error())
		}
	}
}

func CreateRabbitMQConnection(cfg core.RabbitMQConfig, q MQClientConfig, tls *core.TLSConfig, numChannels int, isConsume bool) *RabbitMQConnection {
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
		Channels:       make([]*AmqpChannelWrapper, numChannels),
		publishChannel: make(chan PublishMessage, CHANNEL_BUFFER),
	}

	for i := 0; i < numChannels; i++ {
		ch, err := connection.Channel()
		if err != nil {
			core.Error("Failed to create channel: " + err.Error())
		}

		if isConsume {
			err = ch.Qos(q.ConsumerQos, 0, false)
			if err != nil {
				core.Error("Failed to consumer QoS: " + err.Error())
			}
		}

		wrapper := SetupChannel(ch, q, i, isConsume)
		c.Channels[i] = wrapper

		go c.publishWorker(i)
	}

	return &c
}

func (c *RabbitMQConnection) SendMessage(msg PublishMessage) {
	atomic.AddInt64(&c.publishInProgress, 1)
	c.publishChannel <- msg
}

func (c *RabbitMQConnection) ReceiveMessages(queue string, fn RecvMsgFn) error {
	msgs, err := c.Channels[0].Channel.Consume(
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
		c.Channels[i].Channel.Close()
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

func (r *RealRabbitMQInterface) Connect(cfg core.RabbitMQConfig, q MQClientConfig, tls *core.TLSConfig) {
	r.publish = CreateRabbitMQConnection(cfg, q, tls, 4, false)
	r.consume = CreateRabbitMQConnection(cfg, q, tls, 1, true)
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

func (r *RealRabbitMQInterface) GetConsumerQueueName(id int) string {
	return r.consume.Channels[0].Queues[NotificationQueueId]
}

func (r *RealRabbitMQInterface) WaitForAllMessagesToBeSent() {
	for {
		check := atomic.LoadInt64(&r.publish.publishInProgress)
		if check == 0 {
			break
		}

		time.Sleep(1 * time.Second)
	}
}

var DefaultRabbitMQ = RealRabbitMQInterface{}
