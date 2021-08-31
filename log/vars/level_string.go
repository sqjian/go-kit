// Code generated by "stringer -type=Level -linecomment"; DO NOT EDIT.

package vars

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnknownLevel-0]
	_ = x[None-1]
	_ = x[Debug-2]
	_ = x[Info-3]
	_ = x[Warn-4]
	_ = x[Error-5]
}

const _Level_name = "UnknownLevelNoneDebugInfoWarnError"

var _Level_index = [...]uint8{0, 12, 16, 21, 25, 29, 34}

func (i Level) String() string {
	if i < 0 || i >= Level(len(_Level_index)-1) {
		return "Level(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Level_name[_Level_index[i]:_Level_index[i+1]]
}
