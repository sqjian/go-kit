package zap_test

import (
	"github.com/sqjian/go-kit/log/preset"
	"github.com/sqjian/go-kit/log/provider/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	logger, loggerErr := zap.NewLogger(
		zap.WithFileName("test.log"),
		zap.WithMaxSize(3),
		zap.WithMaxBackups(3),
		zap.WithMaxAge(3),
		zap.WithLevel(preset.Info),
		zap.WithConsole(false),
	)

	if loggerErr != nil {
		t.Fatal(loggerErr)
	}

	{
		t.Log(logger.SetLevelOTF(preset.Warn))
		logger.Debugf("testing infof...")
		logger.Infof("testing Infof...")
		logger.Warnf("testing Warnf...")
		logger.Errorf("testing Errorf...")
	}
	{
		t.Log(logger.SetLevelOTF(preset.Warn))
		logger.Debugf("testing infof...")
		logger.Infof("testing Infof...")
		logger.Warnf("testing Warnf...")
		logger.Errorf("testing Errorf...")
	}
}
