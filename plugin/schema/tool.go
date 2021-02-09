package schema

type Opts interface {
	Set(interface{}, interface{})
	Get(interface{}) (interface{}, bool)
}

type Buffer interface {
	Write([]byte) error
	Read() []byte
}

type PluginTools interface {
	Opts
	Buffer
}
