package hash_test

import (
	"github.com/sqjian/go-kit/hash"
	"testing"
)

func TestGeneratorMd5(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	generator, generatorErr := hash.NewGenerator()
	checkErr(generatorErr)

	calc, calcErr := generator.Calc(hash.MD5, []byte("hello world"))
	checkErr(calcErr)
	t.Logf("calc:%v", calc)
}

func TestGeneratorSha1(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	generator, generatorErr := hash.NewGenerator()
	checkErr(generatorErr)

	calc, calcErr := generator.Calc(hash.SHA1, []byte("hello world"))
	checkErr(calcErr)
	t.Logf("calc:%v", calc)
}
