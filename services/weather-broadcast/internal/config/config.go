package config

import (
	"fmt"
	"strings"
	grpcpkg "weather-forecast/pkg/grpc"
	"weather-forecast/pkg/rabbitmq"

	"github.com/spf13/viper"
)

type Config struct {
	WeatherServiceAddress      string `mapstructure:"WEATHER_SERVICE_ADDRESS"`
	SubscriptionServiceAddress string `mapstructure:"SUBSCRIPTION_SERVICE_ADDRESS"`

	RabbitMQ rabbitmq.Config `mapstructure:",squash"`

	GRPC grpcpkg.Config `mapstructure:",squash"`

	Timezone string `mapstructure:"TIMEZONE"`

	ServiceName       string `mapstructure:"SERVICE_NAME"`
	MetricsServerPort string `mapstructure:"METRICS_SERVER_PORT"`

	LogLevel        string `mapstructure:"LOG_LEVEL"`
	LogSamplingRate int    `mapstructure:"LOG_SAMPLING_RATE"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func validate(config *Config) error {
	if err := config.RabbitMQ.Validate(); err != nil {
		return err
	}

	if err := config.GRPC.Validate(); err != nil {
		return err
	}

	required := map[string]string{
		"WEATHER_SERVICE_ADDRESS":      config.WeatherServiceAddress,
		"SUBSCRIPTION_SERVICE_ADDRESS": config.SubscriptionServiceAddress,
		"TIMEZONE":                     config.Timezone,
		"SERVICE_NAME":                 config.ServiceName,
		"METRICS_SERVER_PORT":          config.MetricsServerPort,
		"LOG_LEVEL":                    config.LogLevel,
	}

	var missing []string
	for name, value := range required {
		if value == "" {
			missing = append(missing, name)
		}
	}

	if config.LogSamplingRate < 1 {
		missing = append(missing, "LOG_SAMPLING_RATE")

	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}

	return nil
}
