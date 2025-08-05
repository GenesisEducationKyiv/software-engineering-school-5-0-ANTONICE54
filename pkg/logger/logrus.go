package logger

import (
	"context"
	"os"
	"weather-forecast/pkg/ctxutil"

	"github.com/sirupsen/logrus"
)

type logrusWrapper struct {
	entry   *logrus.Entry
	sampler Sampler
}

func (l *logrusWrapper) Debugf(format string, args ...interface{}) {
	if l.sampler.ShouldLog() {
		l.entry.Debugf(format, args...)
	}
}

func (l *logrusWrapper) Infof(format string, args ...interface{}) {
	if l.sampler.ShouldLog() {
		l.entry.Infof(format, args...)
	}
}

func (l *logrusWrapper) Warnf(format string, args ...interface{}) {
	if l.sampler.ShouldLog() {
		l.entry.Warnf(format, args...)
	}
}

func (l *logrusWrapper) Errorf(format string, args ...interface{}) {
	if l.sampler.ShouldLog() {
		l.entry.Errorf(format, args...)
	}
}

func (l *logrusWrapper) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusWrapper) WithField(key string, value interface{}) Logger {
	return &logrusWrapper{
		entry:   l.entry.WithField(key, value),
		sampler: l.sampler,
	}
}

func (l *logrusWrapper) WithContext(ctx context.Context) Logger {
	processID := ctxutil.GetCorrelationID(ctx)
	return l.WithField(ctxutil.CorrelationIDKey.String(), processID)
}

func NewLogrus(serviceName, level string, sampler Sampler) *logrusWrapper {
	if sampler == nil {
		sampler = &NoSampler{}
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(toLogrusLevel(level))

	entry := logger.WithField("service", serviceName)
	return &logrusWrapper{
		entry:   entry,
		sampler: sampler,
	}
}
