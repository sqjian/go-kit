package rdb_test

import (
	"context"
	"github.com/sqjian/go-kit/db/rdb"
	"github.com/sqjian/go-kit/log"
	"testing"
	"time"
)

var db *rdb.Rdb

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	_db, dbErr := rdb.NewRdb(
		rdb.Mysql,
		rdb.WithIp("182.92.1.35"),
		rdb.WithPort("3306"),
		rdb.WithUserName("root"),
		rdb.WithPassWord("xylx1.t!@#"),
		rdb.WithMaxLifeTime(time.Second),
		rdb.WithMaxIdleConns(3),
		rdb.WithDbName("test"),
		rdb.WithLogger(log.DebugLogger),
	)
	checkErr(dbErr)

	db = _db
}
func Test_MysqlQuery(t *testing.T) {
	ctx := context.WithValue(context.Background(), "id", 1)

	t.Log(db.Query(ctx, "test", nil))
	t.Log(db.Query(ctx, "test", map[string]interface{}{"t1": 1}))
	t.Log(db.Query(ctx, "test", map[string]interface{}{"t2": 1}))
}

func Test_Insert(t *testing.T) {
	ctx := context.WithValue(context.Background(), "id", 1)

	t.Log(db.Insert(ctx, "test", map[string]interface{}{"t1": 1}))
	t.Log(db.Insert(ctx, "test", map[string]interface{}{"t2": 1}))
}

func Test_Update(t *testing.T) {
	ctx := context.WithValue(context.Background(), "id", 1)

	t.Log(db.Update(ctx, "test", map[string]interface{}{"t1": 5}, map[string]interface{}{"t1": 1}))
}

func Test_Delete(t *testing.T) {
	ctx := context.Background()

	t.Log(db.Delete(ctx, "test", map[string]interface{}{"t2": 1}))
}

func Test_RawSql(t *testing.T) {
	ctx := context.WithValue(context.Background(), "id", 1)

	_, err := db.RawSql(
		ctx,
		"INSERT INTO test (`t1`) VALUES (1)",
		"INSERT INTO test (`t2`) VALUES (1)",
		"INSERT INTO test (`t3`) VALUES (1)",
		"INSERT INTO test (`t4`) VALUES (1)",
		"INSERT INTO test (`t5`) VALUES (1)",
	)
	checkErr(err)
}
