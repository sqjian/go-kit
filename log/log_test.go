package log_test

import (
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/log/vars"
	"testing"
)

func TestLoggerZap(t *testing.T) {
	logger, loggerErr := log.NewLogger(
		log.WithFileName("go-kit.log"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel(vars.Info),
		log.WithLogType(vars.Zap),
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

func TestLoggerDef(t *testing.T) {
	logger, loggerErr := log.NewLogger(
		log.WithFileName("test.log"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel(vars.Info),
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

func TestLoggerDummy(t *testing.T) {
	logger, loggerErr := log.NewLogger(
		log.WithFileName("test.log"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel(vars.Info),
		log.WithLogType(vars.Dummy),
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

func TestDummyLogger(t *testing.T) {

	{
		log.DummyLogger.Debugf("testing infof...")
		log.DummyLogger.Infof("testing Infof...")
		log.DummyLogger.Warnf("testing Warnf...")
		log.DummyLogger.Errorf("testing Errorf...")
	}
}

func TestDebugLogger(t *testing.T) {

	{
		log.DebugLogger.Debugf("testing infof...")
		log.DebugLogger.Infof("testing Infof...")
		log.DebugLogger.Warnf("testing Warnf...")
		log.DebugLogger.Errorf("testing Errorf...")
	}
}
