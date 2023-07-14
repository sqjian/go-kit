package log

func newDummyLogger() *dummyLogger {
	return &dummyLogger{}
}

type dummyLogger struct {
}

func (d dummyLogger) Debugf(template string, args ...interface{}) {
}

func (d dummyLogger) Infof(template string, args ...interface{}) {
}

func (d dummyLogger) Warnf(template string, args ...interface{}) {
}

func (d dummyLogger) Errorf(template string, args ...interface{}) {
}
