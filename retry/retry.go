package retry

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"strings"
	"time"
)

type UserFunc func() error

func Do(userFn UserFunc, opts ...Option) error {
	var n uint

	config := newDefaultRetryConfig()

	for _, opt := range opts {
		opt(config)
	}

	if err := config.context.Err(); err != nil {
		config.logger.Errorf("config.context.Err:%v", err.Error())
		return err
	}

	var errorLog Error

	for n < config.attempts {
		err := userFn()

		if err != nil {
			config.logger.Errorf("call userFn failed,err:%v", err.Error())

			errorLog = append(errorLog, err)

			if !config.retryIf(err) {
				break
			}

			config.onRetry(n, err)

			if n == config.attempts-1 {
				break
			}

			delayTime := config.delayType(n, err, config)

			select {
			case <-time.After(delayTime):
			case <-config.context.Done():
				return append(errorLog, config.context.Err())
			}

		} else {
			return nil
		}
		n++
	}

	return errorLog
}

func newDefaultRetryConfig() *Config {
	return &Config{
		attempts: uint(10),
		logger:   func() log.Log { inst, _ := log.NewLogger(log.WithLevel("dummy")); return inst }(),
		context:  context.Background(),
		onRetry:  func(n uint, err error) {},
		retryIf:  func(err error) bool { return true },
		delayType: func(_ uint, _ error, _ *Config) time.Duration {
			return 100 * time.Millisecond
		},
	}
}

type Error []error

func (e Error) Error() string {
	logWithNumber := make([]string, lenWithoutNil(e))
	for i, l := range e {
		if l != nil {
			logWithNumber[i] = fmt.Sprintf("#%d: %s", i+1, l.Error())
		}
	}
	return fmt.Sprintf("All attempts fail:\n%s", strings.Join(logWithNumber, "\n"))
}

func lenWithoutNil(e Error) (count int) {
	for _, v := range e {
		if v != nil {
			count++
		}
	}

	return
}

func (e Error) WrappedErrors() []error {
	return e
}
