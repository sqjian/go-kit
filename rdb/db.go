package rdb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func newMysqlDb(m *Meta) (*sqlx.DB, error) {
	path := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?collation=utf8mb4_general_ci", m.UserName, m.PassWord, m.IP, m.Port, m.DbName)
	return sqlx.Open("mysql", path)
}

func newSqliteDb(m *Meta) (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", m.DbName+".db")
}

func newDb(dbType Type, dbMeta *Meta) (*sqlx.DB, error) {
	switch dbType {
	case Mysql:
		{
			return newMysqlDb(dbMeta)
		}
	case Sqlite:
		{
			return newSqliteDb(dbMeta)
		}
	default:
		{
			return nil, errWrapper(IllegalParams)
		}
	}
}

func newPlaceHolder(dbType Type) (string, error) {
	switch dbType {
	case Mysql:
		{
			return "?", nil
		}
	case Sqlite:
		{
			return "?", nil
		}
	default:
		{
			return "", errWrapper(IllegalParams)
		}
	}
}
