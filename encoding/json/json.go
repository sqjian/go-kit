package json

import (
	"github.com/buger/jsonparser"
)

func Set(data []byte, setValue []byte, keys ...string) (value []byte, err error) {
	return jsonparser.Set(data, setValue, keys...)
}

func Get(data []byte, keys ...string) (value []byte, dataType jsonparser.ValueType, offset int, err error) {
	return jsonparser.Get(data, keys...)
}
