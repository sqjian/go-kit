package jsonl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/sqjian/go-kit/encoding/json"
	"reflect"
	"strings"
)

func splitJsonL(data []byte, callback func(string)) error {
	var (
		jsonBuffer string
	)

	bracketsCount := 0 // for {}
	squareCount := 0   // for []

	/*
		1、流式去除注释
		2、然后按需截取完整json
		3、解析json
		4、测试用例增加不按行分割的jsonl用例
	*/
	scanner := bufio.NewScanner(bytes.NewReader(Standardize(data) /*这里去除注释的环节不满足流式要求*/))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 {
			continue
		}

		for _, char := range line {
			switch char {
			case '{':
				bracketsCount++
			case '}':
				bracketsCount--
			case '[':
				squareCount++
			case ']':
				squareCount--
			}
		}

		jsonBuffer += line

		if bracketsCount == 0 && squareCount == 0 && len(jsonBuffer) > 0 {
			callback(jsonBuffer)
			jsonBuffer = ""
		}
	}
	return scanner.Err()
}
func Unmarshal(data []byte, ptrToSlice any) error {
	ptr2sl := reflect.TypeOf(ptrToSlice)
	if ptr2sl.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer to slice, got %s", ptr2sl.Kind())
	}

	originalSlice := reflect.Indirect(reflect.ValueOf(ptrToSlice))
	sliceType := originalSlice.Type()
	if sliceType.Kind() != reflect.Slice {
		return fmt.Errorf("expected pointer to slice, got pointer to %s", sliceType.Kind())
	}

	slElem := originalSlice.Type().Elem()

	var decodeErr error
	jsonsErr := splitJsonL(data, func(jsonBuffer string) {
		newObj := reflect.New(slElem).Interface()
		unmarshalErr := json.Unmarshal(Standardize([]byte(jsonBuffer)), newObj)
		if unmarshalErr != nil {
			decodeErr = unmarshalErr
		}
		ptrToNewObj := reflect.Indirect(reflect.ValueOf(newObj))
		originalSlice.Set(reflect.Append(originalSlice, ptrToNewObj))
	})
	if jsonsErr != nil {
		return jsonsErr
	}
	return decodeErr
}

func Marshal(data any) ([]byte, error) {
	originalSlice := reflect.ValueOf(data)
	if originalSlice.Type().Kind() == reflect.Ptr {
		originalSlice = reflect.Indirect(originalSlice)
	}
	if originalSlice.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected slice, got %s", originalSlice.Kind())
	}

	buf := &bytes.Buffer{}

	enc := json.NewEncoder(buf)
	for i := 0; i < originalSlice.Len(); i++ {
		elem := originalSlice.Index(i).Interface()
		err := enc.Encode(elem)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
