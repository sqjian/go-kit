package helper_test

import (
	"github.com/sqjian/go-kit/helper"
	"testing"
)

func TestSet(t *testing.T) {
	set := helper.Make[int]()

	set.Add(1)
	set.Add(3)
	set.Add(5)
	set.Add(7)
	set.Add(1)
	t.Log(set.Contains(2))
}
