package rdb

import (
	"database/sql"
	"fmt"
)

func NewMysqlDb(m *Meta) (*sql.DB, error) {
	path := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", m.UserName, m.PassWord, m.IP, m.Port, m.DbName)
	return sql.Open("sqlite", path)
}
