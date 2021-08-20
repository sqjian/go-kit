package retry_test

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/retry"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	checkErr := func(err error) {
		t.Logf("checkErr:%v", err)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		cnt := 0

		t.Logf("testing group 1...")
		err := retry.Do(
			func() error {
				cnt++
				if cnt < 3 {
					return fmt.Errorf("cnt<3")
				}
				return nil
			},
			retry.WithAttempts(4),
			retry.WithDelay(time.Second),
			retry.WithOnRetry(func(n uint, err error) { t.Logf("NO.%d retry finish,err:%v", n, err) }),
			retry.WithIf(func(err error) bool { t.Logf("ready to retry again,pre err:%v", err); return true }),
			retry.WithContext(func() context.Context { ctx, _ := context.WithTimeout(context.Background(), 3*time.Second); return ctx }()),
		)
		checkErr(err)
	}
}
