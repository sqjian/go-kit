package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/sqjian/go-kit/encoding/jsonc"
	"github.com/sqjian/go-kit/helper"
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
	return helper.RunesToBytes(jsonc.TrimCommentWrapper(helper.BytesToRunes(data)))
}

func Marshal(data any) ([]byte, error) {

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false) // 不对HTML特殊字符进行转义

	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, fmt.Errorf("error encoding JSON:%v", err)
	}

	// 去除末尾默认添加的换行符
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}

	return result, nil
}
