package rdb

import (
	"fmt"
	"github.com/elgris/sqrl"
	"strings"
)

type Instruct struct {
	Sql  string
	Args []interface{}
}

func NewInstruct(sql sqrl.Sqlizer) (*Instruct, error) {
	sqlStr, sqlArgs, sqlErr := sql.ToSql()
	if sqlErr != nil {
		return nil, sqlErr
	}
	return &Instruct{sqlStr, sqlArgs}, nil
}

func instructsToString(instructs []*Instruct) string {
	var rst []string
	for _, instruct := range instructs {
		rst = append(rst, fmt.Sprintf("%v", *instruct))
	}
	return strings.Join(rst, ",")
}

func genQuerySql(table []string, column []string, where map[string]interface{}, group []string, offset, limit uint64) (instruct *Instruct, err error) {

	var sql *sqrl.SelectBuilder

	if len(column) == 0 {
		column = append(column, "*")
	}
	sql = sqrl.Select(column...)

	if len(table) == 0 {
		return nil, fmt.Errorf("empty table")
	}
	sql = sql.From(table...)

	if len(where) != 0 {
		sql = sql.Where(sqrl.Eq(where))
	}

	if len(group) != 0 {
		sql = sql.GroupBy(group...)
	}

	sql = sql.Offset(offset).Limit(limit)

	return NewInstruct(sql)
}

func genInsertSql(table string, data map[string]interface{}) (instruct *Instruct, err error) {

	var sql *sqrl.InsertBuilder

	sql = sqrl.Insert(table)

	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	var columns []string
	var values []interface{}
	for column, value := range data {
		columns = append(columns, column)
		values = append(values, value)
	}
	sql = sql.Columns(columns...)
	sql = sql.Values(values...)

	return NewInstruct(sql)
}

func genDeleteSql(table string, where map[string]interface{}) (instruct *Instruct, err error) {

	var sql *sqrl.DeleteBuilder

	sql = sqrl.Delete(table)

	if len(where) == 0 {
		return nil, fmt.Errorf("empty where")
	}
	sql = sql.Where(where)

	return NewInstruct(sql)
}

func genUpdateSql(table string, data map[string]interface{}, where map[string]interface{}) (instruct *Instruct, err error) {

	var sql *sqrl.UpdateBuilder

	sql = sqrl.Update(table)

	if len(where) == 0 {
		return nil, fmt.Errorf("empty where")
	}
	sql = sql.Where(where)

	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	sql = sql.SetMap(data)

	return NewInstruct(sql)
}
