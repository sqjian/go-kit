// Code generated by "stringer -type=Err -linecomment"; DO NOT EDIT.

package rule

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnknownErrCode-0]
	_ = x[NotFound-1]
}

const _Err_name = "UnknownErrCodeNotFound"

var _Err_index = [...]uint8{0, 14, 22}

func (i Err) String() string {
	if i < 0 || i >= Err(len(_Err_index)-1) {
		return "Err(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Err_name[_Err_index[i]:_Err_index[i+1]]
}
