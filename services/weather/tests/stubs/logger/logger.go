package stub_logger

type StubLogger struct{}

func New() *StubLogger {
	return &StubLogger{}
}

func (l *StubLogger) Debugf(format string, args ...interface{}) {}
func (l *StubLogger) Infof(format string, args ...interface{})  {}
func (l *StubLogger) Warnf(format string, args ...interface{})  {}
func (l *StubLogger) Fatalf(format string, args ...interface{}) {}
