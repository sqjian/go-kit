// Code generated by "stringer -type=LogType -linecomment"; DO NOT EDIT.

package vars

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Zap-1]
	_ = x[Dummy-2]
}

const _LogType_name = "ZapDummy"

var _LogType_index = [...]uint8{0, 3, 8}

func (i LogType) String() string {
	i -= 1
	if i < 0 || i >= LogType(len(_LogType_index)-1) {
		return "LogType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _LogType_name[_LogType_index[i]:_LogType_index[i+1]]
}