package log

func newDummyLogger() *dummyLogger {
	return &dummyLogger{}
}

type dummyLogger struct {
}

func (d dummyLogger) Debugf(template string, args ...any) {
}

func (d dummyLogger) Infof(template string, args ...any) {
}

func (d dummyLogger) Warnf(template string, args ...any) {
}

func (d dummyLogger) Errorf(template string, args ...any) {
}
