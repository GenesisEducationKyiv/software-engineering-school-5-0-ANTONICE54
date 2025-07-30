package config

import (
	"fmt"
	"strings"
	"weather-forecast/pkg/logger"
	"weather-forecast/pkg/rabbitmq"

	"github.com/spf13/viper"
)

type (
	Mailer struct {
		From     string `mapstructure:"MAILER_FROM"`
		Host     string `mapstructure:"MAILER_HOST"`
		Port     string `mapstructure:"MAILER_PORT"`
		Username string `mapstructure:"MAILER_USERNAME"`
		Password string `mapstructure:"MAILER_PASSWORD"`
	}

	Config struct {
		RabbitMQ rabbitmq.Config `mapstructure:",squash"`

		Mailer Mailer `mapstructure:",squash"`

		ServerHost string `mapstructure:"SERVER_HOST"`
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

	if err := config.RabbitMQ.Validate(); err != nil {
		return err
	}

	required := map[string]string{
		"MAILER_FROM":     config.Mailer.From,
		"MAILER_HOST":     config.Mailer.Host,
		"MAILER_PORT":     config.Mailer.Port,
		"MAILER_USERNAME": config.Mailer.Username,
		"MAILER_PASSWORD": config.Mailer.Password,
		"SERVER_HOST":     config.ServerHost,
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
