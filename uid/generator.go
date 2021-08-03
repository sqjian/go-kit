package uid

import (
	"github.com/google/uuid"
	"github.com/sqjian/go-kit/uid/snowflake"
)

type Generator interface {
	NewUuidV1() (string, error)
	NewSnowflake() (string, error)
}

type generator struct {
	MetaData struct {
		nodeId int64 /*current offers for the snowflake */
	}
	Node *snowflake.Node
}

func (g *generator) NewUuidV1() (string, error) {
	i, e := uuid.NewUUID()
	if e != nil {
		return "", e
	}
	return i.String(), nil
}

func (g *generator) NewSnowflake() (string, error) {
	id := g.Node.Generate()
	return id.String(), nil
}

func NewGenerator(opts ...Option) (Generator, error) {

	generatorTmp := new(generator)

	for _, opt := range opts {
		opt.apply(generatorTmp)
	}

	node, nodeErr := snowflake.NewNode(generatorTmp.MetaData.nodeId)
	if nodeErr != nil {
		return nil, nodeErr
	}
	generatorTmp.Node = node

	return generatorTmp, nil
}
