// Code generated by "stringer -type=ErrCode -linecomment"; DO NOT EDIT.

package connection

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[IllegalParams-0]
	_ = x[GetConnTimeout-1]
	_ = x[PoolExhausted-2]
}

const _ErrCode_name = "illegal paramsGet Connection TimeoutPool Was Exhausted"

var _ErrCode_index = [...]uint8{0, 14, 36, 54}

func (i ErrCode) String() string {
	if i < 0 || i >= ErrCode(len(_ErrCode_index)-1) {
		return "ErrCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrCode_name[_ErrCode_index[i]:_ErrCode_index[i+1]]
}