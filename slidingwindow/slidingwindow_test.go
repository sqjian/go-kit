package slidingwindow_test

import (
	"github.com/sqjian/go-kit/slidingwindow"
	"testing"
	"time"
)

type bucket struct {
	collect []any
}

func (b *bucket) empty() bool {
	return len(b.collect) == 0
}

func (b *bucket) Update(a ...any) {
	b.collect = append(b.collect, a...)
}

func (b *bucket) Clear() {
	b.collect = make([]any, 0)
}

func Test_slidingWindow(t *testing.T) {
	genCnt := func(buckets []slidingwindow.Bucket) int {
		var cnt int
		for _, bkt := range buckets {
			if !(bkt.(*bucket).empty()) {
				cnt++
			}
		}
		return cnt
	}

	// 初始化一个bucket为1s，长度3s的窗口
	inst := slidingwindow.NewSlidingWindow(
		time.Second,
		3,
		func() slidingwindow.Bucket {
			return &bucket{}
		},
	)
	{
		inst.Update(1)
		if cnt := genCnt(inst.GetBuckets()); cnt != 1 {
			t.Fatalf("test failed,cnt:%v!=1", cnt)
		}
	}
	{
		t.Logf("sleeping 1s")
		time.Sleep(time.Second)
		inst.Update(2)
		if cnt := genCnt(inst.GetBuckets()); cnt != 2 {
			t.Fatalf("test failed,cnt:%v!=2", cnt)
		}
	}
	{
		t.Logf("sleeping 1s")
		time.Sleep(time.Second)
		inst.Update(3)
		if cnt := genCnt(inst.GetBuckets()); cnt != 3 {
			t.Fatalf("test failed,cnt:%v!=3", cnt)
		}
	}
	{
		t.Logf("sleeping 1s")
		time.Sleep(time.Second)
		inst.Update(4)
		if cnt := genCnt(inst.GetBuckets()); cnt != 3 {
			t.Fatalf("test failed,cnt:%v!=3", cnt)
		}
	}
	{
		t.Logf("sleeping 1s")
		time.Sleep(time.Second)
		if cnt := genCnt(inst.GetBuckets()); cnt != 2 {
			t.Fatalf("test failed,cnt:%v!=2", cnt)
		}
	}
	{
		t.Logf("sleeping 1s")
		time.Sleep(time.Second)
		if cnt := genCnt(inst.GetBuckets()); cnt != 1 {
			t.Fatalf("test failed,cnt:%v!=1", cnt)
		}
	}
	{
		t.Logf("sleeping 1s")
		time.Sleep(time.Second)
		if cnt := genCnt(inst.GetBuckets()); cnt != 0 {
			t.Fatalf("test failed,cnt:%v!=0", cnt)
		}
	}
}
