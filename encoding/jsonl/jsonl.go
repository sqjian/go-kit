package jsonl

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sqjian/go-kit/encoding/jsonc"
	"github.com/sqjian/go-kit/helper"
	"io"
	"reflect"
	"sync"
)

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

	decoder := func(jsonBuffer []byte) error {
		newObj := reflect.New(slElem).Interface()
		var unmarshalErr = json.Unmarshal(helper.RemoveZWNBS(jsonBuffer), newObj)
		if unmarshalErr != nil {
			return unmarshalErr
		}
		ptrToNewObj := reflect.Indirect(reflect.ValueOf(newObj))
		originalSlice.Set(reflect.Append(originalSlice, ptrToNewObj))
		return nil
	}

	return Decode(bytes.NewReader(data), decoder)
}

// Decode data: 数据流，会按json进行拆分并使用decoder处理
func Decode(data io.Reader, decoder func([]byte) error) error {
	var wg sync.WaitGroup

	const chanQueue = 1024

	ctx, cancel := context.WithCancel(context.Background())

	// 读取原始数据
	rawDataChan := func() chan rune {
		ch := make(chan rune, chanQueue)
		wg.Add(1)
		go func() {
			defer func() {
				close(ch)
				wg.Done()
			}()
			buf := bufio.NewReader(data)

			for {
				char, _, err := buf.ReadRune()
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
		}()
		return ch
	}()

	// 去除注释
	trimmedCommentChan := make(chan rune, chanQueue)
	{
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			trimCommentErr := jsonc.TrimComment(ctx, rawDataChan, trimmedCommentChan)
			close(trimmedCommentChan)
			if trimCommentErr != nil {
				cancel()
			}
		}()
	}

	// 切分jsonl
	jsonLSplitChan := make(chan string, chanQueue)
	{
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			splitErr := split(ctx, trimmedCommentChan, jsonLSplitChan)
			close(jsonLSplitChan)
			if splitErr != nil {
				cancel()
			}
		}()
	}

	// 解析json
	var decodeErr error
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		for jsonBuffer := range jsonLSplitChan {
			decodeErr = decoder([]byte(helper.RemoveZWNBS(jsonBuffer)))
			if decodeErr != nil {
				cancel()
			}
		}
	}()

	wg.Wait()

	return decodeErr
}

func split(ctx context.Context, from <-chan rune, to chan<- string) error {

	sendBack := func(str string) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			to <- str
		}
		return nil
	}

	var (
		quote bool // 判断是否在双引号内部
	)

	var (
		jsonBuffer = &bytes.Buffer{}
	)

	bracketsCount := 0 // for {}
	squareCount := 0   // for []
	validCharacters := false

	var preChar rune
	for char := range from {
		jsonBuffer.WriteRune(char)

		if preChar != jsonc.ESCAPE && char == jsonc.QUOTE {
			quote = !quote
		}

		preChar = char

		if quote {
			continue
		}

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
			if err := sendBack(jsonBuffer.String()); err != nil {
				return err
			}
			jsonBuffer.Reset()
			validCharacters = false
		}
	}

	return nil
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
	enc.SetEscapeHTML(false) // 不对HTML特殊字符进行转义
	for i := 0; i < originalSlice.Len(); i++ {
		elem := originalSlice.Index(i).Interface()
		err := enc.Encode(elem)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
