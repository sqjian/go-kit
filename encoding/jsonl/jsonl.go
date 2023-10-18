package jsonl

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sqjian/go-kit/encoding/jsonc"
	"io"
	"reflect"
	"sync"
)

func Unmarshal(data []byte, ptrToSlice any) error {
	return Decode(bytes.NewReader(data), ptrToSlice)
}
func Decode(data io.Reader, ptrToSlice any) error {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	// 读取原始数据
	rawDataChan := func() chan byte {
		ch := make(chan byte)
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := bufio.NewReader(data)

			for {
				char, err := buf.ReadByte()
				if err != nil {
					break
				}
				select {
				case <-ctx.Done():
					return
				default:
					ch <- char
				}
			}
			close(ch)
		}()
		return ch
	}()

	// 去除注释
	commentTrimmedChan := make(chan byte)
	{
		wg.Add(1)
		go func() {
			defer wg.Done()
			jsonc.TrimComment(ctx, rawDataChan, commentTrimmedChan)
		}()
	}

	// 切分jsonl
	jsonLSplitedChan := make(chan string)
	{
		wg.Add(1)
		go func() {
			defer wg.Done()
			split(ctx, commentTrimmedChan, jsonLSplitedChan)
		}()
	}

	// 解析json
	{
		ptr2sl := reflect.TypeOf(ptrToSlice)
		if ptr2sl.Kind() != reflect.Ptr {
			cancel()
			return fmt.Errorf("expected pointer to slice, got %s", ptr2sl.Kind())
		}

		originalSlice := reflect.Indirect(reflect.ValueOf(ptrToSlice))
		sliceType := originalSlice.Type()
		if sliceType.Kind() != reflect.Slice {
			cancel()
			return fmt.Errorf("expected pointer to slice, got pointer to %s", sliceType.Kind())
		}

		slElem := originalSlice.Type().Elem()

		var decodeErr error
		wg.Add(1)
		go func() {
			defer wg.Done()
			for jsonBuffer := range jsonLSplitedChan {
				newObj := reflect.New(slElem).Interface()
				unmarshalErr := json.Unmarshal([]byte(jsonBuffer), newObj)
				if unmarshalErr != nil {
					cancel()
					decodeErr = unmarshalErr
				}
				ptrToNewObj := reflect.Indirect(reflect.ValueOf(newObj))
				originalSlice.Set(reflect.Append(originalSlice, ptrToNewObj))
			}
		}()

		wg.Wait()

		return decodeErr
	}

}

func split(ctx context.Context, from <-chan byte, to chan<- string) {

	var (
		jsonBuffer = &bytes.Buffer{}
	)

	bracketsCount := 0 // for {}
	squareCount := 0   // for []
	validCharacters := false

	for char := range from {
		select {
		case <-ctx.Done():
			return
		default:
		}

		jsonBuffer.WriteByte(char)
		switch char {
		case '{':
			bracketsCount++
			validCharacters = true
		case '}':
			bracketsCount--
			validCharacters = true
		case '[':
			squareCount++
			validCharacters = true
		case ']':
			squareCount--
			validCharacters = true
		}
		if validCharacters && bracketsCount == 0 && squareCount == 0 && jsonBuffer.Len() > 2 {
			to <- jsonBuffer.String()
			jsonBuffer.Reset()
			validCharacters = false
		}
	}

	close(to)
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
