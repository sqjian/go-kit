//go:build dev
// +build dev

package redis

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func Test_Redis(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	inst, instErr := newClusterClient(
		[]string{
			"x.x.x.x:x",
			"x.x.x.x:x",
			"x.x.x.x:x",
		},
		"xxx",
	)
	checkErr(instErr)

	hSetErr := inst.HSet(context.Background(), "xxx", map[string]any{"gender": "男"})
	checkErr(hSetErr)

	get, getErr := inst.HGet(context.Background(), "xxx", "gender")
	checkErr(getErr)

	t.Log(spew.Sdump(get))

}
