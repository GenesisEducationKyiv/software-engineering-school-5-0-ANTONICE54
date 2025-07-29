package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type (
	Retry struct {
		MaxRetries int `mapstructure:"MAILER_MAX_RETRIES"`
		Delay      int `mapstructure:"MAILER_DELAY"`
	}

	Config struct {
		RabbitMQURL string `mapstructure:"RABBIT_MQ_SOURCE"`
		Exchange    string `mapstructure:"EXCHANGE"`

		LogLevel string `mapsturcutre:"LOG_LEVEL"`

		MailerFrom     string `mapstructure:"MAILER_FROM"`
		MailerHost     string `mapstructure:"MAILER_HOST"`
		MailerPort     string `mapstructure:"MAILER_PORT"`
		MailerUsername string `mapstructure:"MAILER_USERNAME"`
		MailerPassword string `mapstructure:"MAILER_PASSWORD"`

		Retry Retry `mapstructure:",squash"`

		ServerHost  string `mapstructure:"SERVER_HOST"`
		ServiceName string `mapstructure:"SERVICE_NAME"`

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
	required := map[string]string{
		"RABBIT_MQ_SOURCE":    config.RabbitMQURL,
		"EXCHANGE":            config.Exchange,
		"MAILER_FROM":         config.MailerFrom,
		"MAILER_HOST":         config.MailerHost,
		"MAILER_PORT":         config.MailerPort,
		"MAILER_USERNAME":     config.MailerUsername,
		"MAILER_PASSWORD":     config.MailerPassword,
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
