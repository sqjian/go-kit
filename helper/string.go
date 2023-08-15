package helper

import "fmt"

func SplitAfter(str string, sep []rune, min int) ([]string, error) {
	if len(sep) == 0 {
		return nil, fmt.Errorf("empty sep")
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
	if cnt != 0 {
		rst = append(rst, string(strRune[len(strRune)-cnt:]))
	}
	return rst, nil
}
