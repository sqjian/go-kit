package log_test

import (
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/log/preset"
	"testing"
)

func TestLogger(t *testing.T) {
	logger, loggerErr := log.NewLogger(
		log.WithFileName("test.log"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel(preset.Info),
		log.WithLogType(preset.Zap),
		log.WithConsole(false),
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
}
