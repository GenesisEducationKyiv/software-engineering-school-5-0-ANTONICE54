package rabbitmq

import (
	"fmt"
	"strings"
)

type Config struct {
	Source     string `mapstructure:"RABBIT_MQ_SOURCE"`
	Retries    int    `mapstructure:"RABBIT_MQ_RETRIES"`
	RetryDelay int    `mapstructure:"RABBIT_MQ_RETRY_DELAY"`
	Exchange   string `mapstructure:"EXCHANGE"`
}

func (c *Config) Validate() error {
	var missing []string

	if c.Source == "" {
		missing = append(missing, "RABBIT_MQ_SOURCE")
	}

	if c.Exchange == "" {
		missing = append(missing, "EXCHANGE")
	}

	if c.Retries < 1 {
		missing = append(missing, "RABBIT_MQ_RETRIES")
	}
	if c.RetryDelay < 1 {
		missing = append(missing, "RABBIT_MQ_RETRY_DELAY")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}

	return nil
}
