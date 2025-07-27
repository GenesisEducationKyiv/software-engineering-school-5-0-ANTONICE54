package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	WeatherServiceAddress      string `mapstructure:"WEATHER_SERVICE_ADDRESS"`
	SubscriptionServiceAddress string `mapstructure:"SUBSCRIPTION_SERVICE_ADDRESS"`

	RabbitMQSource string `mapstructure:"RABBIT_MQ_SOURCE"`
	Exchange       string `mapstructure:"EXCHANGE"`

	Timezone string `mapstructure:"TIMEZONE"`

	ServiceName string `mapstructure:"SERVICE_NAME"`
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
	required := map[string]string{
		"WEATHER_SERVICE_ADDRESS":      config.WeatherServiceAddress,
		"SUBSCRIPTION_SERVICE_ADDRESS": config.SubscriptionServiceAddress,
		"RABBIT_MQ_SOURCE":             config.RabbitMQSource,
		"EXCHANGE":                     config.Exchange,
		"TIMEZONE":                     config.Timezone,
		"SERVICE_NAME":                 config.ServiceName,
	}

	var missing []string
	for name, value := range required {
		if value == "" {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}

	return nil
}
