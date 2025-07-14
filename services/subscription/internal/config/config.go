package config

import (
	"fmt"
	"strings"
	"weather-forecast/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	GRPCPort string `mapstructure:"GRPC_PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`

	RabbitMQSource string `mapstructure:"RABBIT_MQ_SOURCE"`
	Exchange       string `mapstructure:"EXCHANGE"`
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
	required := map[string]string{
		"GRPC_PORT":        config.GRPCPort,
		"DB_HOST":          config.DBHost,
		"DB_USER":          config.DBUser,
		"DB_PASSWORD":      config.DBPassword,
		"DB_NAME":          config.DBName,
		"DB_PORT":          config.DBPort,
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
