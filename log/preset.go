package log

import (
	"fmt"
	"os"
	"strings"
)

type DummyLogger struct{}

func (DummyLogger) Debugf(string, ...interface{}) {}
func (DummyLogger) Infof(string, ...interface{})  {}
func (DummyLogger) Warnf(string, ...interface{})  {}
func (DummyLogger) Errorf(string, ...interface{}) {}

type TerminalLogger struct{}

func Debugf(template string, args ...interface{}) {
	(&TerminalLogger{}).Debugf(template, args...)
}
func (t TerminalLogger) Debugf(template string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, strings.ReplaceAll(template, "\n", "\\n")+"\n", args...)
}
func Infof(template string, args ...interface{}) {
	(&TerminalLogger{}).Debugf(template, args...)
}
func (t TerminalLogger) Infof(template string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, strings.ReplaceAll(template, "\n", "\\n")+"\n", args...)
}
func Warnf(template string, args ...interface{}) {
	(&TerminalLogger{}).Debugf(template, args...)
}
func (t TerminalLogger) Warnf(template string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, strings.ReplaceAll(template, "\n", "\\n")+"\n", args...)
}
func Errorf(template string, args ...interface{}) {
	(&TerminalLogger{}).Debugf(template, args...)
}
func (t TerminalLogger) Errorf(template string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, strings.ReplaceAll(template, "\n", "\\n")+"\n", args...)
}
