package rdb_test

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/rdb"
	"time"
)

var db *rdb.Rdb

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Example_init() {
	_db, dbErr := rdb.NewRdb(
		rdb.Mysql,
		rdb.WithIp("192.168.6.6"),
		rdb.WithPort("3306"),
		rdb.WithUserName("root"),
		rdb.WithPassWord("xylx1.t!@#"),
		rdb.WithMaxLifeTime(time.Second),
		rdb.WithMaxIdleConns(3),
		rdb.WithDbName("test"),
		rdb.WithLogger(func() log.Log { inst, _ := log.NewLogger(log.WithLevel(log.Dummy)); return inst }()),
	)
	checkErr(dbErr)

	db = _db
}

func Example_insert() {
	ctx := context.WithValue(context.Background(), "id", 1)
	fmt.Println(db.Insert(ctx, "test", map[string]interface{}{"age": 1}))
}

func Example_mysqlQuery() {
	ctx := context.WithValue(context.Background(), "id", 1)
	rst, err := db.Query(ctx, []string{"test"}, nil)
	if err != nil {
		panic(err)
	}
	for k, v := range rst {
		fmt.Printf("k:%v,v:%v", k, v)
	}
	fmt.Println(string((rst[0]["create_time"]).([]byte)))
}

func Example_update() {
	ctx := context.WithValue(context.Background(), "id", 1)

	fmt.Println(db.Update(ctx, "test", map[string]interface{}{"age": 5}, map[string]interface{}{"age": 1}))
}

func Example_delete() {
	ctx := context.Background()

	fmt.Println(db.Delete(ctx, "test", map[string]interface{}{"age": 5}))
}
