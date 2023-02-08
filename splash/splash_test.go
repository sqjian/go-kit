package splash_test

import (
	"fmt"
	"github.com/sqjian/go-kit/splash"
	"testing"
)

func TestStringifyMultiLine(t *testing.T) {
	fmt.Println(splash.Stringify("xxx", "xxx", "xxx", "xxx"))
}
