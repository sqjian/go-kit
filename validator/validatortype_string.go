// Code generated by "stringer -type=ValidatorType -linecomment"; DO NOT EDIT.

package validator

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Json-1]
}

const _ValidatorType_name = "Json"

var _ValidatorType_index = [...]uint8{0, 4}

func (i ValidatorType) String() string {
	i -= 1
	if i < 0 || i >= ValidatorType(len(_ValidatorType_index)-1) {
		return "ValidatorType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ValidatorType_name[_ValidatorType_index[i]:_ValidatorType_index[i+1]]
}
