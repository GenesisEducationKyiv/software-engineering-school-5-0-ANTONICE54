package config

import (
	"fmt"
	"strings"
	grpcpkg "weather-forecast/pkg/grpc"
	"weather-forecast/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	WeatherServiceAddress      string `mapstructure:"WEATHER_SERVICE_ADDRESS"`
	SubscriptionServiceAddress string `mapstructure:"SUBSCRIPTION_SERVICE_ADDRESS"`
	ServerPort                 string `mapstructure:"SERVER_PORT"`

	GRPC grpcpkg.Config `mapstructure:",squash"`
}

func Load(log logger.Logger) (*Config, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := validate(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validate(config *Config) error {
	if err := config.GRPC.Validate(); err != nil {
		return err
	}
	required := map[string]string{
		"WEATHER_SERVICE_ADDRESS":      config.WeatherServiceAddress,
		"SUBSCRIPTION_SERVICE_ADDRESS": config.SubscriptionServiceAddress,
		"SERVER_PORT":                  config.ServerPort,
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
