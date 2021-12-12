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
	Consumer
}

type Consumer interface {
	handle(ctx context.Context, deliveries <-chan amqp.Delivery)
	Consume(ctx context.Context) error
}

func (c *scanConsumer) Consume(ctx context.Context) error {
	deliChan, err := c.channel.Consume(false)
	if err != nil {
		return err
	}
	c.handle(ctx, deliChan)
	return nil
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
