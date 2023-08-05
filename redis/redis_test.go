package redis_test

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/sqjian/go-kit/redis"
	"testing"
	"time"
)

func Test_Redis(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	err := redis.Init(
		[]string{
			"node-a.redis-hf04-0ckir5.svc.hfb.ipaas.cn:9000",
			"node-b.redis-hf04-0ckir5.svc.hfb.ipaas.cn:9000",
			"node-c.redis-hf04-0ckir5.svc.hfb.ipaas.cn:9000",
			"node-d.redis-hf04-0ckir5.svc.hfb.ipaas.cn:9000",
			"node-e.redis-hf04-0ckir5.svc.hfb.ipaas.cn:9000",
			"node-f.redis-hf04-0ckir5.svc.hfb.ipaas.cn:9000",
		},
		"xylx1.t!@#",
	)
	checkErr(err)

	err = redis.HSet(context.Background(), "sqjian", "gender", "男")
	checkErr(err)

	time.Sleep(time.Second)

	resp, respErr := redis.HGet(context.Background(), "sqjian", "gender")
	checkErr(respErr)

	t.Log(spew.Sdump(resp))

}
