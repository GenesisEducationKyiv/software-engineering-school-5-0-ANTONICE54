package logger

import (
	"errors"

	"github.com/sirupsen/logrus"
)

const (
	InfoLevel  = "info"
	DebugLevel = "debug"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

func toLogrusLevel(level string) (logrus.Level, error) {
	switch level {
	case DebugLevel:
		return logrus.DebugLevel, nil
	case InfoLevel:
		return logrus.InfoLevel, nil
	case WarnLevel:
		return logrus.WarnLevel, nil
	case ErrorLevel:
		return logrus.ErrorLevel, nil
	default:
		return 0, errors.New("invalid log level")
	}
}
