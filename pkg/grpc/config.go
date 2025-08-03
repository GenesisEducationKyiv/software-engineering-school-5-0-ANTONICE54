package grpc

import (
	"fmt"
	"strings"
)

type Config struct {
	Retries    int `mapstructure:"GRPC_RETRIES"`
	RetryDelay int `mapstructure:"GRPC_RETRY_DELAY"`
}

func (c *Config) Validate() error {
	var missing []string

	if c.Retries < 1 {
		missing = append(missing, "GRPC_RETRIES")
	}
	if c.RetryDelay < 1 {
		missing = append(missing, "GRPC_RETRY_DELAY")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}

	return nil
}
