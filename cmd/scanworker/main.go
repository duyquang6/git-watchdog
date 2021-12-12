package main

import (
	"context"
	"github.com/duyquang6/git-watchdog/internal/configuration"
	"github.com/duyquang6/git-watchdog/internal/consumer"
	"github.com/duyquang6/git-watchdog/internal/core"
	"github.com/duyquang6/git-watchdog/internal/database"
	"github.com/duyquang6/git-watchdog/internal/rabbitmq"
	"github.com/duyquang6/git-watchdog/internal/repository"
	"github.com/duyquang6/git-watchdog/internal/setup"
	"github.com/duyquang6/git-watchdog/pkg/logging"
	"github.com/sethvargo/go-signalcontext"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"time"
)

// main wrap realMain around a graceful shutdown scheme
func main() {
	ctx, done := signalcontext.OnInterrupt()

	logger := logging.NewLoggerFromEnv()
	ctx = logging.WithLogger(ctx, logger)

	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Fatalw("application panic", "panic", r)
		}
	}()

	err := realMain(ctx)
	done()

	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("successful shutdown")
}

func realMain(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	var config configuration.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		logger.Fatal(err)
	}
	ctx = configuration.WithAppConfig(ctx, &config)
	scanConsumer := setupDepScanConsumer(ctx, logger, &config, env.Database(), env.RabbitMQClient())
	return scanConsumer.Consume(ctx)
}

func setupDepScanConsumer(ctx context.Context, logger *zap.SugaredLogger, config *configuration.Config,
	db database.DBFactory, rabbitMQConn *amqp.Connection) *consumer.ScanConsumer {
	scanRepository := repository.NewScanRepository()
	done := make(chan error)
	channel := rabbitmq.NewQueueChannelFromConnection(config.RabbitMQConfig(), rabbitMQConn)

	gitScan := core.NewGitScan(logger, config, config.RuleFilePath)
	scanConsumer := consumer.NewScanConsumer(logger, db,
		scanRepository, gitScan,
		rabbitMQConn,
		channel,
		config.RabbitMQ.ConsumerTag,
		done, config.RabbitMQ.RepositoryScanTaskQueue)

	go func() {
		<-ctx.Done()

		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		shutdownCtx = logging.WithLogger(shutdownCtx, logger)
		defer done()

		if err := db.Close(shutdownCtx); err != nil {
			logger.Error("closed db failed:", err)
		}

		logger.Info("closing amqp conn")
		if err := scanConsumer.Close(shutdownCtx); err != nil {
			logger.Error("closed amqp conn failed:", err)
		}
	}()

	return scanConsumer
}
