package helper

import (
	"reflect"
)

func IsPtr(v any) bool {
	pv := reflect.ValueOf(v)
	return pv.Kind() == reflect.Ptr
}
