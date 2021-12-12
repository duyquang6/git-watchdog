package setup

import (
	"context"
	"fmt"
	"github.com/duyquang6/git-watchdog/internal/rabbitmq"

	"github.com/duyquang6/git-watchdog/internal/database"
	"github.com/duyquang6/git-watchdog/internal/serverenv"
	"github.com/duyquang6/git-watchdog/pkg/logging"
	"github.com/sethvargo/go-envconfig"
)

// DatabaseConfigProvider ensures that the environment config can provide a DB config.
// All binaries in this application connect to the database via the same method.
type DatabaseConfigProvider interface {
	DatabaseConfig() *database.Config
}

// RabbitMQConfigProvider ensures that the environment config can provide RabbitMQ Config.
// All binaries in this application connect to the database via the same method.
type RabbitMQConfigProvider interface {
	RabbitMQConfig() *rabbitmq.Config
}

// Setup runs common initialization code for all servers. See SetupWith.
func Setup(ctx context.Context, config interface{}) (*serverenv.ServerEnv, error) {
	return SetupWith(ctx, config, envconfig.OsLookuper())
}

// SetupWith processes the given configuration using envconfig. It is
// responsible for establishing database connections, resolving secrets, and
// accessing app configs.
func SetupWith(ctx context.Context, config interface{}, l envconfig.Lookuper) (*serverenv.ServerEnv, error) {
	logger := logging.FromContext(ctx)

	// Build a list of options to pass to the server env.
	var serverEnvOpts []serverenv.Option

	// Process first round of environment variables.
	if err := envconfig.ProcessWith(ctx, config, l); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}

	logger.Infow("provided", "config", config)

	// Setup the database connection.
	if provider, ok := config.(DatabaseConfigProvider); ok {
		logger.Info("configuring database")

		dbConfig := provider.DatabaseConfig()
		db, err := database.NewFromEnv(ctx, dbConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database: %v", err)
		}

		// Update serverEnv setup.
		serverEnvOpts = append(serverEnvOpts, serverenv.WithDatabase(db))

		logger.Infow("database", "config", dbConfig)
	}

	// Setup the rabbitmq connection.
	if provider, ok := config.(RabbitMQConfigProvider); ok {
		logger.Info("configuring rabbitmq client")

		rabbitMQConfig := provider.RabbitMQConfig()
		amqpClient, err := rabbitmq.NewFromEnv(ctx, rabbitMQConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to amqp server: %v", err)
		}

		// Update serverEnv setup.
		serverEnvOpts = append(serverEnvOpts, serverenv.WithRabbitMQClient(amqpClient))

		logger.Infow("amqp", "config", rabbitMQConfig)
	}

	return serverenv.New(ctx, serverEnvOpts...), nil
}
