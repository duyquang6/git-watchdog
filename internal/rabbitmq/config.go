package rabbitmq

// Config provide database configuration
type Config struct {
	AMQPConnectionString    string `env:"AMQP_SERVER_URL" json:",omitempty"`
	RepositoryScanTaskQueue string `env:"REPO_SCAN_TASK_QUEUE, default=repo_scan_task" json:",omitempty"`
	ConsumerTag             string `env:"CONSUMER_TAG, default=repo_scan_consumer" json:",omitempty"`
	DLQName                 string `env:"DLQ_NAME, default=dead_letter_queue" json:",omitempty"`
	DLXName                 string `env:"DLX_NAME, default=dlx_exchange" json:",omitempty"`
	DLXKey                  string `env:"DLX_KEY, default=dlx_key" json:",omitempty"`
	DLXKind                 string `env:"DLX_KIND, default=fanout" json:",omitempty"`
}

// RabbitMQConfig get db config
func (c *Config) RabbitMQConfig() *Config {
	return c
}
