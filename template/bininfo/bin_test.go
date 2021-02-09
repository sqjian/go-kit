package bininfo_test

import (
	"fmt"
	"github.com/sqjian/go-kit/template/bininfo"
	"testing"
)

func TestStringifySingleLine(t *testing.T) {
	fmt.Println(bininfo.StringifySingleLine())
}

func TestStringifyMultiLine(t *testing.T) {
	fmt.Println(bininfo.StringifyMultiLine())
}
