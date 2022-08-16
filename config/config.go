package config

// Configuration struct
type Configuration struct {
	Port         string   `mapstructure:"port"`
	DB           Database `mapstructure:"database"`
	RabbitMQ     Broker   `mapstructure:"broker"`
	StooqURL     string   `mapstructure:"stooq-url"`
	JwtSecretKey string   `mapstructure:"SECRET_KEY"`
}

// Database data
type Database struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"DB_PASSWORD"`
}

// Broker data
type Broker struct {
	URL       string `mapstructure:"url"`
	QueueName string `mapstructure:"queue-name"`
}
