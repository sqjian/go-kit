package rdb

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sqjian/go-kit/log"
	"strings"
)

type mysql struct {
	db     *sql.DB
	meta   *Meta
	logger log.API
}

func (m *mysql) init() error {
	//path := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", m.meta.UserName, m.meta.PassWord, m.meta.IP, m.meta.Port, m.meta.DbName)
	//db, dbErr := sql.Open("sqlite", path)
	db, dbErr := sql.Open("sqlite3", "sqlite.db")
	if dbErr != nil {
		return dbErr
	}

	db.SetConnMaxLifetime(m.meta.MaxLifeTime)
	db.SetMaxIdleConns(m.meta.MaxIdleConns)

	pingErr := db.Ping()
	if pingErr != nil {
		return pingErr
	}
	return nil
}

func (m *mysql) Query(sqlInfo string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := m.db.Query(sqlInfo, args...)
	if err != nil {
		return nil, err
	}
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength)
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{}
	for rows.Next() {
		_ = rows.Scan(cache...)

		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{})
		}
		list = append(list, item)
	}
	_ = rows.Close()
	return list, nil

}

func (m *mysql) Delete(table string, where map[string]interface{}) error {
	tx, txErr := m.db.Begin()
	if txErr != nil {
		return txErr
	}
	rawSql := func() string {
		s := fmt.Sprintf("DELETE FROM %v WHERE ", table)
		var w []string
		for k, v := range where {
			w = append(w, fmt.Sprintf("%v = %v", k, v))
		}
		s += strings.Join(w, " AND ")
		return s
	}()
	stmt, stmtErr := tx.Prepare(rawSql)
	if stmtErr != nil {
		return stmtErr
	}
	_, execErr := stmt.Exec()
	if execErr != nil {
		return execErr
	}
	return tx.Commit()
}

func (m *mysql) Insert(table string, data map[string]interface{}) error {
	tx, txErr := m.db.Begin()
	if txErr != nil {
		return txErr
	}
	rawSql := func() string {
		s := fmt.Sprintf("INSERT INTO %v ", table)
		var ks, vs []string
		for k, v := range data {
			ks = append(ks, fmt.Sprintf("`%v`", k))
			vs = append(vs, fmt.Sprintf("%v", v))
		}
		s += fmt.Sprintf("(%v) ", strings.Join(ks, ","))
		s += fmt.Sprintf("VALUES (%v) ", strings.Join(vs, ","))
		return s
	}()
	stmt, stmtErr := tx.Prepare(rawSql)
	if stmtErr != nil {
		return stmtErr
	}
	_, execErr := stmt.Exec()
	if execErr != nil {
		return execErr
	}
	return tx.Commit()
}
