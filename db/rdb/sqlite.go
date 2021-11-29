package rdb

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteDb(m *Meta) (*sql.DB, error) {
	return sql.Open("sqlite3", m.DbName+".db")
}
