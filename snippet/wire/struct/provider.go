package _struct

import "github.com/google/wire"

type Foo int
type Bar int

func ProvideFoo() Foo { return 1 }

func ProvideBar() Bar { return 1 }

type FooBar struct {
	MyFoo Foo
	MyBar Bar
}

var Set = wire.NewSet(
	ProvideFoo,
	ProvideBar,
	wire.Struct(new(FooBar), "MyFoo", "MyBar"))
