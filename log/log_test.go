package log_test

import (
	"fmt"
	"github.com/sqjian/go-kit/log"
	"testing"
)

func TestLoggerZap(t *testing.T) {
	logger, loggerErr := log.NewLogger(
		log.WithFileName("go-kit.log"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel(log.Info),
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
}

func TestLoggerDef(t *testing.T) {
	logger, loggerErr := log.NewLogger(
		log.WithFileName("go-kit.log"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel(log.Info),
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
		log.WithLevel(log.Info),
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

func TestDebugLogger(t *testing.T) {

	{
		fmt.Println("--------------------------------")
		log.DebugLogger.Debugf("testing infof...")
		log.DebugLogger.Infof("testing Infof...")
		log.DebugLogger.Warnf("testing Warnf...")
		log.DebugLogger.Errorf("testing Errorf...")
	}
	{
		fmt.Println("--------------------------------")
		log.Debugf("testing infof...")
		log.Infof("testing Infof...")
		log.Warnf("testing Warnf...")
		log.Errorf("testing Errorf...")
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
	logger, loggerErr := log.NewLogger(
		log.WithFileName("test.log"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel(log.Info),
		log.WithConsole(false),
		log.WithBuilder(func(_ *log.Meta) (log.API, error) {
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
