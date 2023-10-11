package log

func newDummyLogger() *dummyLogger {
	return &dummyLogger{}
}

type dummyLogger struct {
}

func (d dummyLogger) Debugw(_ string, _ ...interface{}) {
}

func (d dummyLogger) Infow(_ string, _ ...interface{}) {
}

func (d dummyLogger) Warnw(_ string, _ ...interface{}) {
}

func (d dummyLogger) Errorw(_ string, _ ...interface{}) {
}

func (d dummyLogger) Debugf(_ string, _ ...any) {
}

func (d dummyLogger) Infof(_ string, _ ...any) {
}

func (d dummyLogger) Warnf(_ string, _ ...any) {
}

func (d dummyLogger) Errorf(_ string, _ ...any) {
}
