package schema

type Opts interface {
	Set(any, any)
	Get(any) (any, bool)
}

type Buffer interface {
	Write([]byte) error
	Read() []byte
}

type PlugTools interface {
	Opts
	Buffer
}
