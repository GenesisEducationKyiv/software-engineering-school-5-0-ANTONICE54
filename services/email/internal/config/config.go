package config

import (
	"fmt"
	"strings"
	"weather-forecast/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	RabbitMQURL string `mapstructure:"RABBIT_MQ_SOURCE"`
	Exchange    string `mapstructure:"EXCHANGE"`

	MailerFrom     string `mapstructure:"MAILER_FROM"`
	MailerHost     string `mapstructure:"MAILER_HOST"`
	MailerPort     string `mapstructure:"MAILER_PORT"`
	MailerUsername string `mapstructure:"MAILER_USERNAME"`
	MailerPassword string `mapstructure:"MAILER_PASSWORD"`

	//
	ServerHost string `mapstructure:"SERVER_HOST"`
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
		"RABBIT_MQ_SOURCE": config.RabbitMQURL,
		"EXCHANGE":         config.Exchange,
		"MAILER_FROM":      config.MailerFrom,
		"MAILER_HOST":      config.MailerHost,
		"MAILER_PORT":      config.MailerPort,
		"MAILER_USERNAME":  config.MailerUsername,
		"MAILER_PASSWORD":  config.MailerPassword,
		"SERVER_HOST":      config.ServerHost,
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
