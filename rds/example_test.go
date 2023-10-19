package rds_test

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/rds"
	"time"
)

var db *rds.Rds

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Example_init() {
	_db, dbErr := rds.NewRds(
		rds.Mysql,
		rds.WithIp("192.168.6.6"),
		rds.WithPort("3306"),
		rds.WithUserName("root"),
		rds.WithPassWord("xylx1.t!@#"),
		rds.WithMaxLifeTime(time.Second),
		rds.WithMaxIdleConns(3),
		rds.WithDbName("test"),
		rds.WithLogger(func() log.Log { inst, _ := log.NewLogger(log.WithLevel("dummy")); return inst }()),
	)
	checkErr(dbErr)

	db = _db
}

func Example_insert() {
	ctx := context.WithValue(context.Background(), "id", 1)
	fmt.Println(db.Insert(ctx, "test", map[string]any{"age": 1}))
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

	fmt.Println(db.Update(ctx, "test", map[string]any{"age": 5}, map[string]any{"age": 1}))
}

func Example_delete() {
	ctx := context.Background()

	fmt.Println(db.Delete(ctx, "test", map[string]any{"age": 5}))
}
