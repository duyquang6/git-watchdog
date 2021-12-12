package configuration

import (
	"context"
	"github.com/duyquang6/git-watchdog/internal/database"
	"github.com/duyquang6/git-watchdog/internal/rabbitmq"
)

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

const (
	appConfigKey = contextKey("appConfig")
)

type Config struct {
	Database     database.Config
	RabbitMQ     rabbitmq.Config
	Port         string `env:"PORT, default=8080"`
	RuleFilePath string `env:"RULE_FILE_PATH" json:",omitempty"`
	TempRootDir  string `env:"TEMP_ROOT_DIR, default=/tmp/git-watchdog"`
}

func (c *Config) DatabaseConfig() *database.Config {
	return &c.Database
}
func (c *Config) RabbitMQConfig() *rabbitmq.Config {
	return &c.RabbitMQ
}

// WithAppConfig creates a new context with the provided AppConfig attached.
func WithAppConfig(ctx context.Context, config *Config) context.Context {
	return context.WithValue(ctx, appConfigKey, config)
}

// FromContext returns the logger stored in the context. If not exist, return default logger
func FromContext(ctx context.Context) *Config {
	if config, ok := ctx.Value(appConfigKey).(*Config); ok {
		return config
	}
	return &Config{}
}
