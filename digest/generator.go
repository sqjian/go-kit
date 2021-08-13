package digest

import (
	"github.com/sqjian/go-kit/digest/provider/md5"
)

//go:generate stringer -type=KeyType  -linecomment
type KeyType int

const (
	MD5 KeyType = iota
)

type Generator interface {
	Calc(KeyType, []byte) (string, error)
}

type generator struct {
}

func (g *generator) Calc(keyType KeyType, data []byte) (string, error) {
	switch keyType {
	case MD5:
		{
			return md5.Calc(data), nil
		}
	}
	return "", GenErr(IllegalKeyType)
}

func NewGenerator(opts ...Option) (Generator, error) {

	generatorInst := new(generator)

	for _, opt := range opts {
		opt.apply(generatorInst)
	}
	return generatorInst, nil
}
