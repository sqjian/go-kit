// Code generated by "stringer -type=Err -linecomment"; DO NOT EDIT.

package err

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnknownCode-0]
	_ = x[Code1-1]
	_ = x[Code2-2]
	_ = x[Code3-3]
}

const _Err_name = "UnknownCodethis is error code1this is error code2this is error code3"

var _Err_index = [...]uint8{0, 11, 30, 49, 68}

func (i Err) String() string {
	if i < 0 || i >= Err(len(_Err_index)-1) {
		return "Err(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Err_name[_Err_index[i]:_Err_index[i+1]]
}
