package json

import (
	"github.com/buger/jsonparser"
	"github.com/tailscale/hujson"
)

func Set(data []byte, setValue []byte, keys ...string) (value []byte, err error) {
	return jsonparser.Set(data, setValue, keys...)
}

func Get(data []byte, keys ...string) (value []byte, dataType jsonparser.ValueType, offset int, err error) {
	return jsonparser.Get(data, keys...)
}

func Standardize(data []byte) ([]byte, error) {
	ast, err := hujson.Parse(data)
	if err != nil {
		return data, err
	}
	ast.Standardize()
	return ast.Pack(), nil
}
