package uid_test

import (
	"github.com/sqjian/toolkit/uid"
	"testing"
)

func TestGenerator_NewUuidV1(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	generator, generatorErr := uid.NewGenerator()
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
	generator, generatorErr := uid.NewGenerator(
		uid.WithNodeId(1),
	)
	checkErr(generatorErr)

	snowflake, snowflakeErr := generator.NewSnowflake()
	checkErr(snowflakeErr)
	t.Logf("snowflake:%v", snowflake)

}
