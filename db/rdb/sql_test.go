package rdb_test

import (
	"context"
	"github.com/sqjian/go-kit/db/rdb"
	"github.com/sqjian/go-kit/log"
	"testing"
	"time"
)

func Test_Query(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	rdb, rdbErr := rdb.NewRdb(
		rdb.Sqlite,
		rdb.WithMaxLifeTime(time.Second),
		rdb.WithMaxIdleConns(10),
		rdb.WithDbName("sqlite"),
		rdb.WithLogger(log.DebugLogger),
	)
	checkErr(rdbErr)

	ctx := context.WithValue(context.Background(), "id", 1)

	t.Log(rdb.Query(ctx, "test", nil))
	t.Log(rdb.Query(ctx, "test", map[string]interface{}{"column_1": 1}))
	t.Log(rdb.Query(ctx, "test", map[string]interface{}{"tah2": 1}))
}

func Test_Insert(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	rdb, rdbErr := rdb.NewRdb(
		rdb.Sqlite,
		rdb.WithMaxLifeTime(time.Second),
		rdb.WithMaxIdleConns(10),
		rdb.WithDbName("sqlite"),
		rdb.WithLogger(log.DebugLogger),
	)
	checkErr(rdbErr)

	ctx := context.WithValue(context.Background(), "id", 1)

	t.Log(rdb.Insert(ctx, "test", map[string]interface{}{"column_1": 1}))
}

func Test_BatchInsert(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	rdb, rdbErr := rdb.NewRdb(
		rdb.Sqlite,
		rdb.WithMaxLifeTime(time.Second),
		rdb.WithMaxIdleConns(10),
		rdb.WithDbName("sqlite"),
		rdb.WithLogger(log.DebugLogger),
	)
	checkErr(rdbErr)

	ctx := context.WithValue(context.Background(), "id", 1)

	t.Log(rdb.Insert(ctx, "test", map[string]interface{}{"column_1": 1}))
}

func Test_Update(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	rdb, rdbErr := rdb.NewRdb(
		rdb.Sqlite,
		rdb.WithMaxLifeTime(time.Second),
		rdb.WithMaxIdleConns(10),
		rdb.WithDbName("sqlite"),
		rdb.WithLogger(log.DebugLogger),
	)
	checkErr(rdbErr)

	ctx := context.WithValue(context.Background(), "id", 1)

	t.Log(rdb.Update(ctx, "test", map[string]interface{}{"column_1": 5}, map[string]interface{}{"column_1": 1}))
}

func Test_Delete(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	rdb, rdbErr := rdb.NewRdb(
		rdb.Sqlite,
		rdb.WithMaxLifeTime(time.Second),
		rdb.WithMaxIdleConns(10),
		rdb.WithDbName("sqlite"),
		rdb.WithLogger(log.DebugLogger),
	)
	checkErr(rdbErr)

	ctx := context.Background()

	t.Log(rdb.Delete(ctx, "test", map[string]interface{}{"column_1": 1}))
}
