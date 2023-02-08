package hash

import (
	"github.com/sqjian/go-kit/hash/md5"
	"github.com/sqjian/go-kit/hash/sha1"
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
	case SHA1:
		{
			return sha1.Calc(data), nil
		}
	}
	return "", ErrWrapper(IllegalKeyType)
}

func NewGenerator(opts ...OptionFunc) (Generator, error) {

	generatorInst := new(generator)

	for _, opt := range opts {
		opt(generatorInst)
	}
	return generatorInst, nil
}
