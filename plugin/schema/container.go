package schema

type Container interface {
	Init() error
	FInit() error
	Add([]Plug) error
	Remove(string) error
	Get(...string) ([]Plug, error)
}
