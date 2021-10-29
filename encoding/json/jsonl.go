package json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

func getOriginalSlice(ptrToSlice interface{}) (slice reflect.Value, err error) {
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

func Decode(r io.Reader, ptrToSlice interface{}) error {
	originalSlice, err := getOriginalSlice(ptrToSlice)
	if err != nil {
		return err
	}

	slElem := originalSlice.Type().Elem()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		newObj := reflect.New(slElem).Interface()
		item := scanner.Bytes()
		err := json.Unmarshal(item, newObj)
		if err != nil {
			return err
		}
		ptrToNewObj := reflect.Indirect(reflect.ValueOf(newObj))
		originalSlice.Set(reflect.Append(originalSlice, ptrToNewObj))
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func Encode(w io.Writer, ptrToSlice interface{}) error {
	originalSlice, err := getOriginalSlice(ptrToSlice)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(w)
	for i := 0; i < originalSlice.Len(); i++ {
		elem := originalSlice.Index(i).Interface()
		err = enc.Encode(elem)
		if err != nil {
			return err
		}
	}
	return nil
}
