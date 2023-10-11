package log_test

import (
	"github.com/sqjian/go-kit/log"
	"testing"
)

func TestLogger(t *testing.T) {
	logger, loggerErr := log.NewLogger(
		log.WithFileName("go-kit.log"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel("info"),
		log.WithConsole(false),
		log.WithCaller(true, 1),
	)

	if loggerErr != nil {
		t.Fatal(loggerErr)
	}

	{
		logger.Debugf("testing infof...")
		logger.Infof("testing Infof...")
		logger.Warnf("testing Warnf...")
		logger.Errorf("testing Errorf...")
	}
	{
		logger.Debugw("haha", "key1", "val1")
	}
}
