package unique_test

import (
	"github.com/sqjian/go-kit/unique"
	"testing"
)

func TestGenerator_NewSnowflake(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	generator, generatorErr := unique.NewGenerator(
		unique.WithSnowflakeNodeId(1),
	)
	checkErr(generatorErr)

	snowflake, snowflakeErr := generator.UniqueKey(unique.Snowflake)
	checkErr(snowflakeErr)
	t.Logf("snowflake:%v", snowflake)

}
func TestGenerator_NewUuidV1(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	generator, generatorErr := unique.NewGenerator()
	checkErr(generatorErr)

	uuidV1, uuidV1Err := generator.UniqueKey(unique.UuidV1)
	checkErr(uuidV1Err)
	t.Logf("uuidV1:%v", uuidV1)
}
