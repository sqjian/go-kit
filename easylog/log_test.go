package easylog_test

import (
	"fmt"
	"github.com/sqjian/go-kit/easylog"
	"testing"
)

func TestLoggerZap(t *testing.T) {
	logger, loggerErr := easylog.NewLogger(
		easylog.WithFileName("go-kit.easylog"),
		easylog.WithMaxSize(3),
		easylog.WithMaxBackups(3),
		easylog.WithMaxAge(3),
		easylog.WithLevel(easylog.Info),
		easylog.WithConsole(false),
		easylog.WithCaller(true, 1),
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
	logger, loggerErr := easylog.NewLogger(
		easylog.WithFileName("go-kit.easylog"),
		easylog.WithMaxSize(3),
		easylog.WithMaxBackups(3),
		easylog.WithMaxAge(3),
		easylog.WithLevel(easylog.Info),
		easylog.WithConsole(false),
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
	logger, loggerErr := easylog.NewLogger(
		easylog.WithFileName("test.easylog"),
		easylog.WithMaxSize(3),
		easylog.WithMaxBackups(3),
		easylog.WithMaxAge(3),
		easylog.WithLevel(easylog.Info),
		easylog.WithConsole(false),
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

func TestDebugLogger(t *testing.T) {

	{
		fmt.Println("--------------------------------")
		easylog.DebugLogger.Debugf("testing infof...")
		easylog.DebugLogger.Infof("testing Infof...")
		easylog.DebugLogger.Warnf("testing Warnf...")
		easylog.DebugLogger.Errorf("testing Errorf...")
	}
	{
		fmt.Println("--------------------------------")
		easylog.Debugf("testing infof...")
		easylog.Infof("testing Infof...")
		easylog.Warnf("testing Warnf...")
		easylog.Errorf("testing Errorf...")
	}
}

type Builder struct{}

func (b Builder) Debugf(template string, args ...interface{}) {
	fmt.Println("implement me")
}

func (b Builder) Infof(template string, args ...interface{}) {
	fmt.Println("implement me")
}

func (b Builder) Warnf(template string, args ...interface{}) {
	fmt.Println("implement me")
}

func (b Builder) Errorf(template string, args ...interface{}) {
	fmt.Println("implement me")
}

func TestBuilder(t *testing.T) {
	logger, loggerErr := easylog.NewLogger(
		easylog.WithFileName("test.easylog"),
		easylog.WithMaxSize(3),
		easylog.WithMaxBackups(3),
		easylog.WithMaxAge(3),
		easylog.WithLevel(easylog.Info),
		easylog.WithConsole(false),
		easylog.WithBuilder(func(_ *easylog.Meta) (easylog.API, error) {
			return &Builder{}, nil
		}),
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
