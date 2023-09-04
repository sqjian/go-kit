package log

import (
	"testing"
)

func TestLogger(t *testing.T) {
	logger := newZapLogger(&config{
		FileName:   "go-kit.log",
		MaxSize:    3,
		MaxBackups: 3,
		MaxAge:     3,
		Level:      Debug,
		Console:    true,
	})

	{
		t.Log(logger.SetLevelOTF(Warn))
		logger.Debugf("testing infof...")
		logger.Infof("testing Infof...")
		logger.Warnf("testing Warnf...")
		logger.Errorf("testing Errorf...")
		logger.Errorw("testing Errorf...", "key1", "val1", "key2", "val2", "key3", "val3")
	}
	{
		t.Log(logger.SetLevelOTF(Warn))
		logger.Debugf("testing infof...")
		logger.Infof("testing Infof...")
		logger.Warnf("testing Warnf...")
		logger.Errorf("testing Errorf...")
		logger.Errorw("testing Errorf...", "key1", "val1", "key2", "val2", "key3", "val3")
	}
}

func TestLoggerCaller(t *testing.T) {
	logger := newZapLogger(&config{
		FileName:   "go-kit.log",
		MaxSize:    3,
		MaxBackups: 3,
		MaxAge:     3,
		Level:      Debug,
		Console:    true,
		Caller:     true,
		CallerSkip: 4,
	})

	{
		logger.Debugf("testing infof...")
		logger.Infof("testing Infof...")
		logger.Warnf("testing Warnf...")
		logger.Errorf("testing Errorf...")
	}
}
