package json

import (
	"github.com/buger/jsonparser"
	"github.com/sqjian/go-kit/encoding/jsonc"
)

func Set(data []byte, setValue []byte, keys ...string) (value []byte, err error) {
	return jsonparser.Set(data, setValue, keys...)
}

func Get(data []byte, keys ...string) (value []byte, dataType jsonparser.ValueType, offset int, err error) {
	return jsonparser.Get(data, keys...)
}

func Standardize(data []byte) []byte {
	return jsonc.TrimCommentWrapper(data)
}
