package helper

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

func SplitAfter(str string, sep []rune, min int) ([]string, error) {
	if len(sep) == 0 {
		return nil, fmt.Errorf("illegal sep:%v", spew.Sdump(sep))
	}
	if min <= 0 {
		return nil, fmt.Errorf("illegal min:%v", spew.Sdump(min))
	}

	strRune := []rune(str)

	sepSet := func() Set[rune] {
		sepSetTmp := MakeSet[rune]()
		for _, sepItem := range sep {
			sepSetTmp.Add(sepItem)
		}
		return sepSetTmp
	}()

	cnt := 0
	var rst []string
	for ix, item := range strRune {
		cnt++
		itemStr := string(item)
		_ = itemStr
		if sepSet.Contains(item) {
			if cnt >= min {
				rst = append(rst, string(strRune[ix-cnt+1:ix+1]))
				cnt = 0
			}
		}
	}
	if cnt > 0 {
		return append(rst, string(strRune[len(strRune)-cnt:])), nil
	}
	return rst, nil
}
