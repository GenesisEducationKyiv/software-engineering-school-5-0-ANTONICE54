package config

import (
	"fmt"
	"strings"
	"weather-forecast/pkg/logger"

	"github.com/spf13/viper"
)

type (
	DB struct {
		Host     string `mapstructure:"DB_HOST"`
		User     string `mapstructure:"DB_USER"`
		Password string `mapstructure:"DB_PASSWORD"`
		Name     string `mapstructure:"DB_NAME"`
		Port     string `mapstructure:"DB_PORT"`
	}

	Config struct {
		GRPCPort string `mapstructure:"GRPC_PORT"`

		DB DB `mapstructure:",squash"`

		RabbitMQSource string `mapstructure:"RABBIT_MQ_SOURCE"`
		Exchange       string `mapstructure:"EXCHANGE"`
	}
)

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
	required := map[string]string{
		"GRPC_PORT":        config.GRPCPort,
		"DB_HOST":          config.DB.Host,
		"DB_USER":          config.DB.User,
		"DB_PASSWORD":      config.DB.Password,
		"DB_NAME":          config.DB.Name,
		"DB_PORT":          config.DB.Port,
		"RABBIT_MQ_SOURCE": config.RabbitMQSource,
		"EXCHANGE":         config.Exchange,
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
