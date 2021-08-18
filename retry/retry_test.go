package retry_test

import (
	"context"
	"github.com/sqjian/go-kit/retry"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		t.Logf("testing group 1...")
		err := retry.Do(
			func() error { return nil },
			retry.WithAttempts(3),
			retry.WithDelay(time.Second),
			retry.WithOnRetry(func(n uint, err error) { t.Logf("NO.%d retry finish,err:%v", n, err) }),
			retry.WithIf(func(err error) bool { t.Logf("ready to retry again,pre err:%v", err); return true }),
			retry.WithContext(func() context.Context { ctx, _ := context.WithTimeout(context.Background(), 3*time.Second); return ctx }()),
		)
		checkErr(err)
	}
}
