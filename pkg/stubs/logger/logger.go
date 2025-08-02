package stub_logger

import (
	"context"
	"weather-forecast/pkg/logger"
)

type StubLogger struct{}

func New() *StubLogger {
	return &StubLogger{}
}

func (l *StubLogger) Debugf(format string, args ...interface{}) {}
func (l *StubLogger) Infof(format string, args ...interface{})  {}
func (l *StubLogger) Warnf(format string, args ...interface{})  {}
func (l *StubLogger) Fatalf(format string, args ...interface{}) {}
func (l *StubLogger) Errorf(format string, args ...interface{}) {}
func (l *StubLogger) WithField(key string, value interface{}) logger.Logger {
	return l
}
func (l *StubLogger) WithContext(ctx context.Context) logger.Logger {
	return l
}
