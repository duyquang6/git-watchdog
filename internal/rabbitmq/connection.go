package rabbitmq

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
)

// NewFromEnv sets up the rabbitmq connections using the configuration in the
// process's environment variables. This should be called just once per server
// instance.
func NewFromEnv(ctx context.Context, cfg *Config) (*amqp.Connection, error) {
	connectRabbitMQ, err := amqp.Dial(cfg.AMQPConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connect rabbitmq: %v", err)
	}

	return connectRabbitMQ, nil
}
