package log

type Logger interface {
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})

	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})

	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

type logger struct {
	logType  KeyType
	metaData struct {
		FileName   string /*日志的名字*/
		MaxSize    int    /*日志大小，单位MB*/
		MaxBackups int    /*日志备份个数*/
		MaxAge     int    /*日志备份时间，单位Day*/
		Level      Level  /*日志级别，可选：none、debug、info、warn、error*/
		Console    bool   /*是否向控制台输出*/
	}
}

func NewLogger(opts ...Option) (Generator, error) {

	loggerInst := new(logger)

	for _, opt := range opts {
		opt.apply(loggerInst)
	}

	switch loggerInst.metaData. {

	}

	return loggerInst, nil
}
