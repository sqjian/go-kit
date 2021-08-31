package tag_test

import (
	"fmt"
	"github.com/sqjian/go-kit/template/tag"
	"testing"
)

func TestStringifySingleLine(t *testing.T) {
	fmt.Println(tag.StringifySingleLine())
}

func TestStringifyMultiLine(t *testing.T) {
	fmt.Println(tag.StringifyMultiLine())
}
