package log

import (
	"testing"
)

func TestLogger(t *testing.T) {
	logger, loggerErr := newZapLogger(&Config{
		FileName:   "go-kit.easylog",
		MaxSize:    3,
		MaxBackups: 3,
		MaxAge:     3,
		Level:      Debug,
		Console:    true,
	})

	if loggerErr != nil {
		t.Fatal(loggerErr)
	}

	{
		t.Log(logger.SetLevelOTF(Warn))
		logger.Debugf("testing infof...")
		logger.Infof("testing Infof...")
		logger.Warnf("testing Warnf...")
		logger.Errorf("testing Errorf...")
	}
	{
		t.Log(logger.SetLevelOTF(Warn))
		logger.Debugf("testing infof...")
		logger.Infof("testing Infof...")
		logger.Warnf("testing Warnf...")
		logger.Errorf("testing Errorf...")
	}
}

func TestLoggerCaller(t *testing.T) {
	logger, loggerErr := newZapLogger(&Config{
		FileName:   "go-kit.easylog",
		MaxSize:    3,
		MaxBackups: 3,
		MaxAge:     3,
		Level:      Debug,
		Console:    true,
		Caller:     true,
		CallerSkip: 4,
	})

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
