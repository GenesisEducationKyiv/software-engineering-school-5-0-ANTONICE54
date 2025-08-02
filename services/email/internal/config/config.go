package config

import (
	"fmt"
	"strings"
	"weather-forecast/pkg/rabbitmq"

	"github.com/spf13/viper"
)

type (
	Retry struct {
		MaxRetries int `mapstructure:"MAILER_MAX_RETRIES"`
		Delay      int `mapstructure:"MAILER_DELAY"`
	}

	Mailer struct {
		From     string `mapstructure:"MAILER_FROM"`
		Host     string `mapstructure:"MAILER_HOST"`
		Port     string `mapstructure:"MAILER_PORT"`
		Username string `mapstructure:"MAILER_USERNAME"`
		Password string `mapstructure:"MAILER_PASSWORD"`
	}

	Config struct {
		RabbitMQ rabbitmq.Config `mapstructure:",squash"`

		Retry Retry `mapstructure:",squash"`

		Mailer Mailer `mapstructure:",squash"`

		ServerHost  string `mapstructure:"SERVER_HOST"`
		ServiceName string `mapstructure:"SERVICE_NAME"`

		LogLevel          string `mapstructure:"LOG_LEVEL"`
		MetricsServerPort string `mapstructure:"METRICS_SERVER_PORT"`
		LogSamplingRate   int    `mapstructure:"LOG_SAMPLING_RATE"`
	}
)

func Load() (*Config, error) {
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
		"MAILER_FROM":         config.Mailer.From,
		"MAILER_HOST":         config.Mailer.Host,
		"MAILER_PORT":         config.Mailer.Port,
		"MAILER_USERNAME":     config.Mailer.Username,
		"MAILER_PASSWORD":     config.Mailer.Password,
		"SERVER_HOST":         config.ServerHost,
		"SERVICE_NAME":        config.ServiceName,
		"METRICS_SERVER_PORT": config.MetricsServerPort,
		"LOG_LEVEL":           config.LogLevel,
	}

	var missing []string
	for name, value := range required {
		if value == "" {
			missing = append(missing, name)
		}
	}

	if config.Retry.Delay == 0 {
		missing = append(missing, "MAILER_DELAY")
	}
	if config.Retry.MaxRetries == 0 {
		missing = append(missing, "MAILER_MAX_RETRIES")
	}
	if config.LogSamplingRate == 0 {
		missing = append(missing, "LOG_SAMPLING_RATE")

	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}

	return nil
}
