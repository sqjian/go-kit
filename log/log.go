package log

type Config struct {
	FileName   string /*日志的名字*/
	MaxSize    int    /*日志大小，单位MB*/
	MaxBackups int    /*日志备份个数*/
	MaxAge     int    /*日志备份时间，单位Day*/
	Level      Level  /*日志级别，可选：none、debug、info、warn、error*/
	Console    bool   /*是否向控制台输出*/
	Caller     bool
	CallerSkip int

	builder func(*Config) (API, error)
}

func newDefaultConfig() *Config {
	return &Config{
		CallerSkip: 1,
	}
}

func NewLogger(opts ...OptionFunc) (API, error) {

	config := newDefaultConfig()

	for _, opt := range opts {
		opt(config)
	}

	switch {
	case len(config.FileName) == 0:
		return nil, ErrWrapper(IllegalParams)
	case config.MaxSize == 0:
		return nil, ErrWrapper(IllegalParams)
	case config.MaxAge == 0:
		return nil, ErrWrapper(IllegalParams)
	case config.Level == UnknownLevel:
		return nil, ErrWrapper(IllegalParams)
	}

	if config.builder != nil {
		return config.builder(config)
	}
	return newZapLogger(config)
}
