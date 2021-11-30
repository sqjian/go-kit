package rdb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func newMysqlDb(m *Meta) (*sql.DB, error) {
	path := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", m.UserName, m.PassWord, m.IP, m.Port, m.DbName)
	return sql.Open("mysql", path)
}

func newSqliteDb(m *Meta) (*sql.DB, error) {
	return sql.Open("sqlite3", m.DbName+".db")
}

func newDb(dbType Type, meta *Meta) (*sql.DB, error) {
	switch dbType {
	case Mysql:
		{
			return newMysqlDb(meta)
		}
	case Sqlite:
		{
			return newSqliteDb(meta)
		}
	default:
		{
			return nil, errWrapper(IllegalParams)
		}
	}
}
