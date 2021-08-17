package retry_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/sqjian/go-kit/retry"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDoAllFailed(t *testing.T) {
	var retrySum uint
	err := retry.Do(
		func() error { return errors.New("test") },
		retry.OnRetry(func(n uint, err error) { retrySum += n }),
		retry.Delay(time.Nanosecond),
	)
	assert.Error(t, err)

	expectedErrorFormat := `All attempts fail:
#1: test
#2: test
#3: test
#4: test
#5: test
#6: test
#7: test
#8: test
#9: test
#10: test`
	assert.Equal(t, expectedErrorFormat, err.Error(), "retry error format")
	assert.Equal(t, uint(45), retrySum, "right count of retry")
}

func TestDoFirstOk(t *testing.T) {
	var retrySum uint
	err := retry.Do(
		func() error { return nil },
		retry.OnRetry(func(n uint, err error) { retrySum += n }),
	)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), retrySum, "no retry")

}

func TestRetryIf(t *testing.T) {
	var retryCount uint
	err := retry.Do(
		func() error {
			if retryCount >= 2 {
				return errors.New("special")
			} else {
				return errors.New("test")
			}
		},
		retry.OnRetry(func(n uint, err error) { retryCount++ }),
		retry.If(func(err error) bool {
			return err.Error() != "special"
		}),
		retry.Delay(time.Nanosecond),
	)
	assert.Error(t, err)

	expectedErrorFormat := `All attempts fail:
#1: test
#2: test
#3: special`
	assert.Equal(t, expectedErrorFormat, err.Error(), "retry error format")
	assert.Equal(t, uint(2), retryCount, "right count of retry")

}

func TestDefaultSleep(t *testing.T) {
	start := time.Now()
	err := retry.Do(
		func() error { return errors.New("test") },
		retry.Attempts(3),
	)
	dur := time.Since(start)
	assert.Error(t, err)
	assert.True(t, dur > 300*time.Millisecond, "3 times default retry is longer then 300ms")
}

func TestFixedSleep(t *testing.T) {
	start := time.Now()
	err := retry.Do(
		func() error { return errors.New("test") },
		retry.Attempts(3),
		retry.DelayType(retry.FixedDelay),
	)
	dur := time.Since(start)
	assert.Error(t, err)
	assert.True(t, dur < 500*time.Millisecond, "3 times default retry is shorter then 500ms")
}

func TestLastErrorOnly(t *testing.T) {
	var retrySum uint
	err := retry.Do(
		func() error { return fmt.Errorf("%d", retrySum) },
		retry.OnRetry(func(n uint, err error) { retrySum += 1 }),
		retry.Delay(time.Nanosecond),
		retry.LastErrorOnly(true),
	)
	assert.Error(t, err)
	assert.Equal(t, "9", err.Error())
}

func TestUnrecoverableError(t *testing.T) {
	attempts := 0
	expectedErr := errors.New("error")
	err := retry.Do(
		func() error {
			attempts++
			return retry.Unrecoverable(expectedErr)
		},
		retry.Attempts(2),
		retry.LastErrorOnly(true),
	)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 1, attempts, "unrecoverable error broke the loop")
}

func TestCombineFixedDelays(t *testing.T) {
	start := time.Now()
	err := retry.Do(
		func() error { return errors.New("test") },
		retry.Attempts(3),
		retry.DelayType(retry.CombineDelay(retry.FixedDelay, retry.FixedDelay)),
	)
	dur := time.Since(start)
	assert.Error(t, err)
	assert.True(t, dur > 400*time.Millisecond, "3 times combined, fixed retry is longer then 400ms")
	assert.True(t, dur < 500*time.Millisecond, "3 times combined, fixed retry is shorter then 500ms")
}

func TestRandomDelay(t *testing.T) {
	start := time.Now()
	err := retry.Do(
		func() error { return errors.New("test") },
		retry.Attempts(3),
		retry.DelayType(retry.RandomDelay),
		retry.MaxJitter(50*time.Millisecond),
	)
	dur := time.Since(start)
	assert.Error(t, err)
	assert.True(t, dur > 2*time.Millisecond, "3 times random retry is longer then 2ms")
	assert.True(t, dur < 100*time.Millisecond, "3 times random retry is shorter then 100ms")
}

func TestMaxDelay(t *testing.T) {
	start := time.Now()
	err := retry.Do(
		func() error { return errors.New("test") },
		retry.Attempts(5),
		retry.Delay(10*time.Millisecond),
		retry.MaxDelay(50*time.Millisecond),
	)
	dur := time.Since(start)
	assert.Error(t, err)
	assert.True(t, dur > 170*time.Millisecond, "5 times with maximum delay retry is longer than 170ms")
	assert.True(t, dur < 200*time.Millisecond, "5 times with maximum delay retry is shorter than 200ms")
}

func TestCombineDelay(t *testing.T) {
	f := func(d time.Duration) retry.DelayTypeFunc {
		return func(_ uint, _ error, _ *retry.Config) time.Duration {
			return d
		}
	}
	const max = time.Duration(1<<63 - 1)
	for _, c := range []struct {
		label    string
		delays   []time.Duration
		expected time.Duration
	}{
		{
			label: "empty",
		},
		{
			label: "single",
			delays: []time.Duration{
				time.Second,
			},
			expected: time.Second,
		},
		{
			label: "negative",
			delays: []time.Duration{
				time.Second,
				-time.Millisecond,
			},
			expected: time.Second - time.Millisecond,
		},
		{
			label: "overflow",
			delays: []time.Duration{
				max,
				time.Second,
				time.Millisecond,
			},
			expected: max,
		},
	} {
		t.Run(
			c.label,
			func(t *testing.T) {
				funcs := make([]retry.DelayTypeFunc, len(c.delays))
				for i, d := range c.delays {
					funcs[i] = f(d)
				}
				actual := retry.CombineDelay(funcs...)(0, nil, nil)
				assert.Equal(t, c.expected, actual, "delay duration mismatch")
			},
		)
	}
}

func TestContext(t *testing.T) {
	const defaultDelay = 100 * time.Millisecond
	t.Run("cancel before", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		retrySum := 0
		start := time.Now()
		err := retry.Do(
			func() error { return errors.New("test") },
			retry.OnRetry(func(n uint, err error) { retrySum += 1 }),
			retry.Context(ctx),
		)
		dur := time.Since(start)
		assert.Error(t, err)
		assert.True(t, dur < defaultDelay, "immediately cancellation")
		assert.Equal(t, 0, retrySum, "called at most once")
	})

	t.Run("cancel in retry progress", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		retrySum := 0
		err := retry.Do(
			func() error { return errors.New("test") },
			retry.OnRetry(func(n uint, err error) {
				retrySum += 1
				if retrySum > 1 {
					cancel()
				}
			}),
			retry.Context(ctx),
		)
		assert.Error(t, err)
		assert.Equal(t, 2, retrySum, "called at most once")
	})
}
