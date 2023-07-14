package log

type config struct {
	FileName   string /*日志的名字*/
	MaxSize    int    /*日志大小，单位MB*/
	MaxBackups int    /*日志备份个数*/
	MaxAge     int    /*日志备份时间，单位Day*/
	Level      Level  /*日志级别，可选：none、debug、info、warn、error*/
	Console    bool   /*是否向控制台输出*/
	Caller     bool
	CallerSkip int
}

func newDefaultConfig() *config {
	return &config{
		CallerSkip: 1,
	}
}

func NewLogger(opts ...OptionFunc) (Log, error) {

	config := newDefaultConfig()

	for _, opt := range opts {
		opt(config)
	}

	if config.Level == Dummy {
		return newDummyLogger(), nil
	}

	if len(config.FileName) == 0 ||
		config.MaxSize == 0 ||
		config.MaxAge == 0 ||
		config.Level == UnknownLevel {
		return nil, ErrWrapper(IllegalParams)
	}
	return newZapLogger(config), nil
}
