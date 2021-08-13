package bininfo_test

import (
	"fmt"
	"github.com/sqjian/go-kit/bininfo"
	"testing"
)

func TestStringifySingleLine(t *testing.T) {
	fmt.Println(bininfo.StringifySingleLine())
}

func TestStringifyMultiLine(t *testing.T) {
	fmt.Println(bininfo.StringifyMultiLine())
}