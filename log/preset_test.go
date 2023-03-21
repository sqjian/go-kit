package log

import "testing"

func TestTerminalLogger_Debugf(t1 *testing.T) {
	Debugf("print1\n")
	Debugf("print2\n")
}
