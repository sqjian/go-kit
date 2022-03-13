package log

type Meta struct {
	FileName   string /*日志的名字*/
	MaxSize    int    /*日志大小，单位MB*/
	MaxBackups int    /*日志备份个数*/
	MaxAge     int    /*日志备份时间，单位Day*/
	Level      Level  /*日志级别，可选：none、debug、info、warn、error*/
	Console    bool   /*是否向控制台输出*/

	builder func(*Meta) (API, error)
}

func newDefaultMeta() *Meta {
	return &Meta{}
}

func NewLogger(opts ...Option) (API, error) {

	meta := newDefaultMeta()

	for _, opt := range opts {
		opt.apply(meta)
	}

	switch {
	case len(meta.FileName) == 0:
		return nil, ErrWrapper(IllegalParams)
	case meta.MaxSize == 0:
		return nil, ErrWrapper(IllegalParams)
	case meta.MaxAge == 0:
		return nil, ErrWrapper(IllegalParams)
	case meta.Level == UnknownLevel:
		return nil, ErrWrapper(IllegalParams)
	}

	if meta.builder != nil {
		return meta.builder(meta)
	}
	return newZapLogger(meta)
}
