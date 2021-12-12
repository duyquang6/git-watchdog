package rabbitmq

import (
	"github.com/streadway/amqp"
)

var _ IChannel = (*QueueChannel)(nil)

type QueueChannel struct {
	channel   *amqp.Channel
	queueName string
}

func NewQueueChannel(queueName string, channel *amqp.Channel) *QueueChannel {
	return &QueueChannel{
		channel:   channel,
		queueName: queueName,
	}
}

func NewQueueChannelFromConnection(config *Config, amqpClient *amqp.Connection) *QueueChannel {
	channelRabbitMQ, err := amqpClient.Channel()
	if err != nil {
		panic(err)
	}

	args := map[string]interface{}{
		"x-dead-letter-exchange":    config.DLXName,
		"x-dead-letter-routing-key": config.DLXKey,
	}

	initDLQ(channelRabbitMQ, config)
	_, err = channelRabbitMQ.QueueDeclare(
		config.RepositoryScanTaskQueue, // queue name
		true,
		false,
		false,
		false,
		args,
	)

	if err != nil {
		panic(err)
	}
	return &QueueChannel{
		channel:   channelRabbitMQ,
		queueName: config.RepositoryScanTaskQueue,
	}
}

func initDLQ(channelRabbitMQ *amqp.Channel, config *Config) {
	err := channelRabbitMQ.ExchangeDeclare(
		config.DLXName,
		config.DLXKind,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		panic(err)
	}
	_, err = channelRabbitMQ.QueueDeclare(config.DLQName, true,
		false,
		false,
		false,
		nil)
	if err != nil {
		panic(err)
	}
	err = channelRabbitMQ.QueueBind(config.DLQName, config.DLXKind, config.DLXName, false, nil)
	if err != nil {
		panic(err)
	}
}

type IChannel interface {
	PublishToQueue(msg amqp.Publishing) error
	Consume(autoAck bool) (<-chan amqp.Delivery, error)
	Close() error
	Cancel(tag string) error
}

func (s *QueueChannel) Close() error {
	return s.channel.Close()
}

func (s *QueueChannel) Cancel(tag string) error {
	return s.channel.Cancel(tag, false)
}

func (s *QueueChannel) PublishToQueue(msg amqp.Publishing) error {
	return s.channel.Publish("", s.queueName, false, false, msg)
}

func (s *QueueChannel) Consume(autoAck bool) (<-chan amqp.Delivery, error) {
	return s.channel.Consume(s.queueName, "", autoAck, false, false, false, nil)
}
