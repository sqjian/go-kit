package rdb

type Option interface {
	apply(*Meta)
}

type optionFunc func(*Meta)

func (f optionFunc) apply(log *Meta) {
	f(log)
}

func WithUserName(UserName string) Option {
	return optionFunc(func(m *Meta) {
		m.UserName = UserName
	})
}

func WithPassWord(PassWord string) Option {
	return optionFunc(func(m *Meta) {
		m.PassWord = PassWord
	})
}

func WithIp(ip string) Option {
	return optionFunc(func(m *Meta) {
		m.IP = ip
	})
}

func WithPort(port string) Option {
	return optionFunc(func(m *Meta) {
		m.Port = port
	})
}

func WithDbName(dbName string) Option {
	return optionFunc(func(m *Meta) {
		m.DbName = dbName
	})
}
