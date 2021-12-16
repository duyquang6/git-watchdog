package consumer

import (
	"context"
	"fmt"
	"github.com/duyquang6/git-watchdog/internal/rabbitmq"
	"github.com/duyquang6/git-watchdog/pkg/logging"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type BaseConsumer struct {
	conn      *amqp.Connection
	channel   rabbitmq.IChannel
	tag       string
	done      chan error
	queueName string
	logger    *zap.SugaredLogger
	consumer  Consumer
}

type Consumer interface {
	processingMessage(ctx context.Context, message rabbitmq.IDelivery) error
}

func (c *BaseConsumer) Consume(ctx context.Context) error {
	deliChan, err := c.channel.Consume(false)
	if err != nil {
		return err
	}
	c.handle(ctx, deliChan)
	return nil
}

func (c *BaseConsumer) handle(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for delivery := range deliveries {
		d := rabbitmq.NewDelivery(delivery)
		err := c.consumer.processingMessage(ctx, d)
		if err != nil {
			c.logger.Errorf("failed to handle message, "+
				"move message tag %d to dead letter queue", d.DeliveryTag())
			if d.Nack(false, false) != nil {
				c.logger.Error("cannot nack message with delivery_tag:", d.DeliveryTag())
			}
		} else {
			if d.Ack(false) != nil {
				c.logger.Error("cannot ack message with delivery_tag:", d.DeliveryTag())
			}
		}
	}

	c.logger.Info("handle: deliveries channel closed")
	c.done <- nil
}

func (c *BaseConsumer) Close(ctx context.Context) error {
	logger := logging.FromContext(ctx)
	defer func() {
		select {
		// wait for handle() to exit
		case <-c.done:
		default:
		}
	}()
	if err := c.channel.Cancel(c.tag); err != nil {

		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}
	logger.Info("AMQP shutdown OK")

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.done:
		return nil
	}
}
