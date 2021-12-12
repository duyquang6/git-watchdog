package serverenv

import (
	"context"
	"github.com/streadway/amqp"

	"github.com/duyquang6/git-watchdog/internal/database"
)

// ServerEnv provide overall connection, environment app
type ServerEnv struct {
	database       database.DBFactory
	rabbitMQClient *amqp.Connection
}

// Option defines function types to modify the ServerEnv on creation.
type Option func(*ServerEnv) *ServerEnv

// New creates a new ServerEnv with the requested options.
func New(ctx context.Context, opts ...Option) *ServerEnv {
	env := &ServerEnv{}

	for _, f := range opts {
		env = f(env)
	}

	return env
}

// WithDatabase attached a database to the environment.
func WithDatabase(db database.DBFactory) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.database = db
		return s
	}
}

// Database get database
func (s *ServerEnv) Database() database.DBFactory {
	return s.database
}

// WithRabbitMQClient attached a database to the environment.
func WithRabbitMQClient(cli *amqp.Connection) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.rabbitMQClient = cli
		return s
	}
}

// RabbitMQClient get database
func (s *ServerEnv) RabbitMQClient() *amqp.Connection {
	return s.rabbitMQClient
}

// Close shuts down the server env, closing database connections, etc.
func (s *ServerEnv) Close(ctx context.Context) error {
	if s == nil {
		return nil
	}

	if s.database != nil {
		s.database.Close(ctx)
	}

	return nil
}
