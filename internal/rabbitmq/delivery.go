package rabbitmq

import "github.com/streadway/amqp"

type Delivery struct {
	delivery amqp.Delivery
}

func NewDelivery(delivery amqp.Delivery) Delivery {
	return Delivery{delivery}
}

type IDelivery interface {
	Ack(multiple bool) error
	Nack(multiple, requeue bool) error
	Body() []byte
	DeliveryTag() uint64
	ConsumerTag() string
	ContentType() string
}

func (d Delivery) Ack(multiple bool) error {
	return d.delivery.Ack(multiple)
}

func (d Delivery) Nack(multiple, requeue bool) error {
	return d.delivery.Nack(multiple, requeue)
}

func (d Delivery) Body() []byte {
	return d.delivery.Body
}

func (d Delivery) DeliveryTag() uint64 {
	return d.delivery.DeliveryTag
}

func (d Delivery) ConsumerTag() string {
	return d.delivery.ConsumerTag
}

func (d Delivery) ContentType() string {
	return d.delivery.ContentType
}
