// Code generated by "stringer -type=ErrCode -linecomment"; DO NOT EDIT.

package preset

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnknownErrCode-0]
	_ = x[IllegalParams-1]
	_ = x[IllegalKeyType-2]
}

const _ErrCode_name = "UnknownErrCodeIllegalParamsIllegalKeyType"

var _ErrCode_index = [...]uint8{0, 14, 27, 41}

func (i ErrCode) String() string {
	if i < 0 || i >= ErrCode(len(_ErrCode_index)-1) {
		return "ErrCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ErrCode_name[_ErrCode_index[i]:_ErrCode_index[i+1]]
}
