package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	GRPCPort          string `mapstructure:"GRPC_PORT"`
	MetricsServerPort string `mapstructure:"METRICS_SERVER_PORT"`

	RedisSource string `mapstructure:"REDIS_SOURCE"`

	WeatherAPIURL string `mapstructure:"WEATHER_API_URL"`
	WeatherAPIKey string `mapstructure:"WEATHER_API_KEY"`

	OpenWeatherURL string `mapstructure:"OPEN_WEATHER_URL"`
	OpenWeatherKey string `mapstructure:"OPEN_WEATHER_KEY"`

	LogFilePath string `mapstructure:"LOG_FILE_PATH"`
	ServiceName string `mapstructure:"SERVICE_NAME"`

	LogLevel        string `mapstructure:"LOG_LEVEL"`
	LogSamplingRate int    `mapstructure:"LOG_SAMPLING_RATE"`
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
		"GRPC_PORT":           config.GRPCPort,
		"METRICS_SERVER_PORT": config.MetricsServerPort,
		"REDIS_SOURCE":        config.RedisSource,
		"WEATHER_API_URL":     config.WeatherAPIURL,
		"WEATHER_API_KEY":     config.WeatherAPIKey,
		"OPEN_WEATHER_URL":    config.OpenWeatherURL,
		"OPEN_WEATHER_KEY":    config.OpenWeatherKey,
		"LOG_FILE_PATH":       config.LogFilePath,
		"SERVICE_NAME":        config.ServiceName,
		"LOG_LEVEL":           config.LogLevel,
	}

	var missing []string
	for name, value := range required {
		if value == "" {
			missing = append(missing, name)
		}
	}
	if config.LogSamplingRate == 0 {
		missing = append(missing, "LOG_SAMPLING_RATE")

	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}

	return nil
}
