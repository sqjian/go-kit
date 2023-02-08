package rdb_test

import (
	"context"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/rdb"
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
		rdb.WithIp("192.168.6.6"),
		rdb.WithPort("3306"),
		rdb.WithUserName("root"),
		rdb.WithPassWord("xylx1.t!@#"),
		rdb.WithMaxLifeTime(time.Second),
		rdb.WithMaxIdleConns(3),
		rdb.WithDbName("test"),
		rdb.WithLogger(log.TerminalLogger{}),
	)
	checkErr(dbErr)

	db = _db
}

func Test_Insert(t *testing.T) {
	ctx := context.WithValue(context.Background(), "id", 1)
	log.TerminalLogger{}.Debugf("begin to insert data.")
	t.Log(db.Insert(ctx, "test", map[string]interface{}{"age": 1}))
}

func Test_MysqlQuery(t *testing.T) {
	ctx := context.WithValue(context.Background(), "id", 1)
	rst, err := db.Query(ctx, []string{"test"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range rst {
		t.Logf("k:%v,v:%v", k, v)
	}
	t.Log(string((rst[0]["create_time"]).([]byte)))
}

func Test_Update(t *testing.T) {
	ctx := context.WithValue(context.Background(), "id", 1)

	t.Log(db.Update(ctx, "test", map[string]interface{}{"age": 5}, map[string]interface{}{"age": 1}))
}

func Test_Delete(t *testing.T) {
	ctx := context.Background()

	t.Log(db.Delete(ctx, "test", map[string]interface{}{"age": 5}))
}
