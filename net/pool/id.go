package pool

import (
	"golang.org/x/net/context"
	"strconv"
	"sync/atomic"
)

var uniqueIdSeed uint64

func genUniqueId(ctx context.Context) string {
	id := ctx.Value("id")
	if idStr, idOk := id.(string); idOk {
		return idStr
	}
	return "->" + strconv.FormatUint(atomic.AddUint64(&uniqueIdSeed, 1), 10)
}
