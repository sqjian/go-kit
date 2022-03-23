package unique

import (
	"github.com/google/uuid"
	"github.com/sqjian/go-kit/unique/snowflake"
)

//go:generate stringer -type=KeyType  -linecomment
type KeyType int

const (
	Snowflake KeyType = iota
	UuidV1
)

type Generator interface {
	UniqueKey(KeyType) (string, error)
}

type generator struct {
	snowflake struct {
		id   int64 /*current offers for the snowflake */
		Node *snowflake.Node
	}
}

func (g *generator) UniqueKey(keyType KeyType) (string, error) {
	switch keyType {
	case Snowflake:
		{
			id := g.snowflake.Node.Generate()
			return id.String(), nil
		}
	case UuidV1:
		{
			i, e := uuid.NewUUID()
			if e != nil {
				return "", e
			}
			return i.String(), nil
		}
	}
	return "", ErrWrapper(IllegalKeyType)
}

func NewGenerator(opts ...Option) (Generator, error) {

	generatorInst := new(generator)

	for _, opt := range opts {
		opt.apply(generatorInst)
	}

	node, nodeErr := snowflake.NewNode(generatorInst.snowflake.id)
	if nodeErr != nil {
		return nil, nodeErr
	}
	generatorInst.snowflake.Node = node

	return generatorInst, nil
}
