package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type logrusWrapper struct {
	entry *logrus.Entry
}

func (l *logrusWrapper) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}
func (l *logrusWrapper) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}
func (l *logrusWrapper) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}
func (l *logrusWrapper) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}
func (l *logrusWrapper) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logrusWrapper) WithField(key string, value interface{}) Logger {
	return &logrusWrapper{entry: l.entry.WithField(key, value)}
}

func NewLogrus(serviceName string) *logrusWrapper {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)

	entry := logger.WithField("service", serviceName)
	return &logrusWrapper{entry: entry}
}
