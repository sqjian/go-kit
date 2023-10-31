package uid_test

//func TestGenerator_NewSnowflake(t *testing.T) {
//	checkErr := func(err error) {
//		if err != nil {
//			t.Fatal(err)
//		}
//	}
//	generator, generatorErr := uid.NewGenerator(
//		uid.Snowflake,
//		uid.WithSnowflakeNodeId(1),
//	)
//	checkErr(generatorErr)
//
//	snowflake, snowflakeErr := generator.Gen()
//	checkErr(snowflakeErr)
//	t.Logf("snowflake:%v", snowflake)
//
//}
//func TestGenerator_NewUuidV1(t *testing.T) {
//	checkErr := func(err error) {
//		if err != nil {
//			t.Fatal(err)
//		}
//	}
//	generator, generatorErr := uid.NewGenerator(uid.Snowflake)
//	checkErr(generatorErr)
//
//	uuidV1, uuidV1Err := generator.Gen()
//	checkErr(uuidV1Err)
//	t.Logf("uuidV1:%v", uuidV1)
//}
