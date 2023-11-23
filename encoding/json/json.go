package json

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/sqjian/go-kit/encoding/jsonc"
	"strings"
)

func Set(data []byte, setValue []byte, keys ...string) (value []byte, err error) {
	return jsonparser.Set(data, setValue, keys...)
}

func MustGet(data []byte, keys ...string) []byte {
	val, _, _, err := Get(data, keys...)
	if err != nil {
		panic(fmt.Errorf("failed:%w,info:can't reach %v", err, strings.Join(keys, "->")))
	}
	return val
}
func Get(data []byte, keys ...string) (value []byte, dataType jsonparser.ValueType, offset int, err error) {
	return jsonparser.Get(data, keys...)
}

func Standardize(data []byte) []byte {
	return jsonc.TrimCommentWrapper(data)
}
