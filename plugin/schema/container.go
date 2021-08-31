package schema

type Container interface {
	Init() error
	FInit() error
	Add([]Plugin) error
	Remove(string) error
	Get(...string) ([]Plugin, error)
}
