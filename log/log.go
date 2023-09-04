package log

type Log interface {
	Debugf(template string, args ...any)
	Debugw(msg string, keysAndValues ...interface{})
	Infof(template string, args ...any)
	Infow(msg string, keysAndValues ...interface{})
	Warnf(template string, args ...any)
	Warnw(msg string, keysAndValues ...interface{})
	Errorf(template string, args ...any)
	Errorw(msg string, keysAndValues ...interface{})
}
