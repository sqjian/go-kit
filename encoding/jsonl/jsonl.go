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

func getOriginalSlice(ptrToSlice any) (slice reflect.Value, err error) {
	ptr2sl := reflect.TypeOf(ptrToSlice)
	if ptr2sl.Kind() != reflect.Ptr {
		return reflect.ValueOf(nil), fmt.Errorf("expected pointer to slice, got %s", ptr2sl.Kind())
	}

	originalSlice := reflect.Indirect(reflect.ValueOf(ptrToSlice))
	sliceType := originalSlice.Type()
	if sliceType.Kind() != reflect.Slice {
		return reflect.ValueOf(nil), fmt.Errorf("expected pointer to slice, got pointer to %s", sliceType.Kind())
	}
	return originalSlice, nil
}

func Unmarshal(data []byte, ptrToSlice any) error {
	originalSlice, err := getOriginalSlice(ptrToSlice)
	if err != nil {
		return err
	}

	var jsonBuffer string
	bracketsCount := 0 // for {}
	squareCount := 0   // for []

	slElem := originalSlice.Type().Elem()
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		newObj := reflect.New(slElem).Interface()
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
			standardizeData, standardizeDataErr := Standardize([]byte(jsonBuffer))
			if standardizeDataErr != nil {
				return standardizeDataErr
			}
			unmarshalErr := json.Unmarshal(standardizeData, newObj)
			if unmarshalErr != nil {
				return unmarshalErr
			}
			ptrToNewObj := reflect.Indirect(reflect.ValueOf(newObj))
			originalSlice.Set(reflect.Append(originalSlice, ptrToNewObj))
			jsonBuffer = ""
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func Marshal(ptrToSlice any) ([]byte, error) {
	originalSlice, err := getOriginalSlice(ptrToSlice)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	enc := json.NewEncoder(buf)
	for i := 0; i < originalSlice.Len(); i++ {
		elem := originalSlice.Index(i).Interface()
		err = enc.Encode(elem)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
