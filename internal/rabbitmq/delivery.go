package rabbitmq

import "github.com/streadway/amqp"

type delivery struct {
	delivery amqp.Delivery
}

func NewDelivery(deli amqp.Delivery) delivery {
	return delivery{deli}
}

type IDelivery interface {
	Ack(multiple bool) error
	Nack(multiple, requeue bool) error
	Body() []byte
	DeliveryTag() uint64
	ConsumerTag() string
	ContentType() string
}

func (d delivery) Ack(multiple bool) error {
	return d.delivery.Ack(multiple)
}

func (d delivery) Nack(multiple, requeue bool) error {
	return d.delivery.Nack(multiple, requeue)
}

func (d delivery) Body() []byte {
	return d.delivery.Body
}

func (d delivery) DeliveryTag() uint64 {
	return d.delivery.DeliveryTag
}

func (d delivery) ConsumerTag() string {
	return d.delivery.ConsumerTag
}

func (d delivery) ContentType() string {
	return d.delivery.ContentType
}
