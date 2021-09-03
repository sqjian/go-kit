package digest_test

import (
	"github.com/sqjian/go-kit/digest"
	"testing"
)

func TestGeneratorMd5(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	generator, generatorErr := digest.NewGenerator()
	checkErr(generatorErr)

	calc, calcErr := generator.Calc(digest.MD5, []byte("hello world"))
	checkErr(calcErr)
	t.Logf("calc:%v", calc)
}

func TestGeneratorSha1(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	generator, generatorErr := digest.NewGenerator()
	checkErr(generatorErr)

	calc, calcErr := generator.Calc(digest.SHA1, []byte("hello world"))
	checkErr(calcErr)
	t.Logf("calc:%v", calc)
}
