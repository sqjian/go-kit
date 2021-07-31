package id_test

import (
	"github.com/sqjian/toolkit/id"
	"testing"
)

func TestGenerator_NewUuidV1(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	generator, generatorErr := id.NewGenerator()
	checkErr(generatorErr)

	uuidV1, uuidV1Err := generator.NewUuidV1()
	checkErr(uuidV1Err)
	t.Logf("uuidV1:%v", uuidV1)

}

func TestGenerator_NewSnowflake(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	generator, generatorErr := id.NewGenerator(
		id.WithNodeId(1),
	)
	checkErr(generatorErr)

	snowflake, snowflakeErr := generator.NewSnowflake()
	checkErr(snowflakeErr)
	t.Logf("snowflake:%v", snowflake)

}
