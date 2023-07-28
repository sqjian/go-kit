package binding

import "github.com/google/wire"

type Footer interface {
	Foo() string
}

type MyFooter string

func (b *MyFooter) Foo() string {
	return string(*b)
}

func provideMyFooter() *MyFooter {
	b := new(MyFooter)
	*b = "Hello, World!"
	return b
}

type Bar string

func provideBar(f Footer) string {
	return f.Foo()
}

var Set = wire.NewSet(
	provideMyFooter,
	wire.Bind(new(Footer), new(*MyFooter)),
	provideBar)
