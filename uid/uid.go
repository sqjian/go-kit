package uid

import (
	"github.com/google/uuid"
	"github.com/sqjian/go-kit/uid/snowflake"
)

//go:generate stringer -type=KeyType  -linecomment
type UidType int

const (
	_ UidType = iota
	Snowflake
	UuidV1
)

type Uid interface {
	Gen() (string, error)
}

type generator struct {
	uidType UidType

	snowflake struct {
		id   int64 /*current offers for the snowflake */
		Node *snowflake.Node
	}
}

func (g *generator) Gen() (string, error) {
	switch g.uidType {
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

func NewGenerator(keyType UidType, opts ...OptionFunc) (Uid, error) {

	inst := &generator{uidType: keyType}
	for _, opt := range opts {
		opt(inst)
	}

	node, nodeErr := snowflake.NewNode(inst.snowflake.id)
	if nodeErr != nil {
		return nil, nodeErr
	}
	inst.snowflake.Node = node

	return inst, nil
}
