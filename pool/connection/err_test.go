package connection_test

import (
	"github.com/sqjian/go-kit/pool/connection"
	"testing"
)

func TestErrWrapper(t *testing.T) {
	t.Log(connection.ErrWrapper(connection.PoolExhausted, "x.x.x.x:x"))
}
