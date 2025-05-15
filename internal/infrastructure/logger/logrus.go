package logger

import "github.com/sirupsen/logrus"

func NewLogrus() *logrus.Logger {
	logger := logrus.New()

	logger.SetLevel(logrus.DebugLevel)

	return logger
}
