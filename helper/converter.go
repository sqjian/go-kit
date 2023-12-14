package helper

import (
	"bytes"
	"strings"
)

func BytesToRunes(b []byte) []rune {
	return []rune(string(b))
}

func RunesToBytes(runes []rune) []byte {
	return []byte(string(runes))
}

func RemoveZWNBS[T string | []byte](input T) T {
	zwnbs := []byte("\ufeff")
	switch v := any(input).(type) {
	case string:
		return any(strings.ReplaceAll(v, string(zwnbs), "")).(T)
	case []byte:
		return any(bytes.ReplaceAll(v, zwnbs, nil)).(T)
	}
	return input
}
